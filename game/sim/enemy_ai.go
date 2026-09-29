package sim

import (
	"fmt"
	"math"

	"github.com/co0p/tankism-revival/game/ecs"
	ecscomponents "github.com/co0p/tankism-revival/game/ecs/components"
	"github.com/co0p/tankism-revival/lib/mathutil"
)

const (
	enemyAIDecisionIntervalSeconds = 0.35
	enemyAIPathCellSize            = 32.0
	enemyAIPathPadding             = 2.0
	enemyAIWaypointArrivalDistance = 10.0

	enemyAIFireAlignmentThreshold = 0.22
	enemyAIAggressiveFireRange    = 420.0
	enemyAICautiousFireRange      = 380.0
	enemyAIFleeingFireRange       = 180.0

	enemyAICautiousFireChance = 0.25
	enemyAIFleeingFireChance  = 0.10

	enemyAICautiousCoverSearchRadius = 8
	enemyAIFleeSearchRadius          = 10
	enemyAIFriendInfluenceRadius     = 220.0
	enemyAIMaxFriendInfluenceCount   = 3

	enemyAIFleeCoverBonus = 60000.0
)

type behaviorWeights struct {
	aggressive float64
	cautious   float64
	fleeing    float64
}

type navCell struct {
	col int
	row int
}

type navGrid struct {
	cols     int
	rows     int
	cellSize float64
	blocked  []bool
}

type neighborStep struct {
	dCol int
	dRow int
	cost int
}

var navNeighborSteps = []neighborStep{
	{dCol: 0, dRow: -1, cost: 10},
	{dCol: 1, dRow: -1, cost: 14},
	{dCol: 1, dRow: 0, cost: 10},
	{dCol: 1, dRow: 1, cost: 14},
	{dCol: 0, dRow: 1, cost: 10},
	{dCol: -1, dRow: 1, cost: 14},
	{dCol: -1, dRow: 0, cost: 10},
	{dCol: -1, dRow: -1, cost: 14},
}

func NewEnemyAI(entity ecs.Entity) ecscomponents.EnemyAI {
	rngState := seedForEnemyEntity(entity)
	favored := sampleFavoredBehavior(&rngState)

	return ecscomponents.EnemyAI{
		Behavior:                favored,
		FavoredBehavior:         favored,
		DecisionCooldownSeconds: 0,
		RNGState:                rngState,
		Path:                    nil,
		PathIndex:               0,
	}
}

func UpdateEnemyTankIntents(world *ecs.World, target, mapEntity ecs.Entity, stopDistance, deltaSeconds float64) error {
	for _, enemy := range world.EnemyAITanks() {
		if err := UpdateEnemyTankIntent(world, enemy.Entity, target, mapEntity, stopDistance, deltaSeconds); err != nil {
			return err
		}
	}

	return nil
}

func UpdateEnemyTankIntent(world *ecs.World, enemy, target, mapEntity ecs.Entity, stopDistance, deltaSeconds float64) error {
	enemyTransform, ok := world.Transform(enemy)
	if !ok {
		return fmt.Errorf("enemy transform missing")
	}

	targetTransform, ok := world.Transform(target)
	if !ok {
		return fmt.Errorf("target transform missing")
	}

	enemyAIState := normalizedEnemyAIState(world, enemy)

	decisionTick := false
	enemyAIState.DecisionCooldownSeconds -= deltaSeconds
	if enemyAIState.DecisionCooldownSeconds <= 0 {
		decisionTick = true
		enemyAIState.DecisionCooldownSeconds += enemyAIDecisionIntervalSeconds
		for enemyAIState.DecisionCooldownSeconds <= 0 {
			enemyAIState.DecisionCooldownSeconds += enemyAIDecisionIntervalSeconds
		}
	}

	healthRatio := enemyHealthRatio(world, enemy)
	cooldownReady := enemyCooldownReady(world, enemy)
	nearbyFriends := countNearbyEnemyFriends(world, enemy, enemyTransform)
	if decisionTick {
		weights := markovTransitionWeights(enemyAIState.Behavior, enemyAIState.FavoredBehavior, healthRatio, cooldownReady, nearbyFriends)
		enemyAIState.Behavior = sampleBehavior(weights, &enemyAIState.RNGState)
	}

	obstacles := world.SolidObstacleColliderBodies()
	clearShot := lineOfSightClear(enemyTransform.X, enemyTransform.Y, targetTransform.X, targetTransform.Y, obstacles)

	if decisionTick {
		mapBounds, err := mapBoundsFromEntity(world, mapEntity)
		if err != nil {
			return err
		}

		enemyCollider := enemyColliderRadius(world, enemy)
		nav := buildNavGrid(mapBounds, obstacles, enemyCollider)
		startCell := nav.worldToCell(enemyTransform.X, enemyTransform.Y)

		targetCell := plannedTargetCellForBehavior(nav, enemyAIState.Behavior, startCell, targetTransform, obstacles)
		path := findPathAStar(nav, startCell, targetCell)
		enemyAIState.Path = enemyAIState.Path[:0]
		enemyAIState.PathIndex = 0
		if len(path) > 1 {
			enemyAIState.Path = make([]ecscomponents.GridWaypoint, 0, len(path)-1)
			for _, step := range path[1:] {
				enemyAIState.Path = append(enemyAIState.Path, ecscomponents.GridWaypoint{Col: step.col, Row: step.row})
			}
		}
	}

	steerX, steerY, hasSteer := nextSteeringWaypoint(&enemyAIState, enemyTransform)
	if !hasSteer {
		switch enemyAIState.Behavior {
		case ecscomponents.EnemyBehaviorAggressive:
			steerX = targetTransform.X
			steerY = targetTransform.Y
			hasSteer = true
		case ecscomponents.EnemyBehaviorCautious:
			if clearShot {
				steerX, steerY = awayPoint(enemyTransform, targetTransform, 120)
				hasSteer = true
			}
		case ecscomponents.EnemyBehaviorFleeing:
			steerX, steerY = awayPoint(enemyTransform, targetTransform, 160)
			hasSteer = true
		}
	}

	intent := ecscomponents.TankControlIntent{}
	turnDeltaToPlayer, distanceToPlayer := turnAndDistanceToPoint(enemyTransform, targetTransform.X, targetTransform.Y)

	if hasSteer {
		turnDeltaToSteer, _ := turnAndDistanceToPoint(enemyTransform, steerX, steerY)
		intent.Turn = mathutil.ClampUnit(turnDeltaToSteer / (math.Pi / 4))

		throttle := behaviorThrottle(enemyAIState.Behavior, clearShot, distanceToPlayer, stopDistance)
		if throttle > 0 && math.Abs(turnDeltaToSteer) > (math.Pi*0.75) {
			throttle *= 0.25
		}
		intent.Throttle = mathutil.ClampUnit(throttle)
	}

	alignedToPlayer := math.Abs(turnDeltaToPlayer) <= enemyAIFireAlignmentThreshold

	stateFireRollPassed := false
	if decisionTick {
		switch enemyAIState.Behavior {
		case ecscomponents.EnemyBehaviorCautious:
			stateFireRollPassed = nextRandomUnit(&enemyAIState.RNGState) < enemyAICautiousFireChance
		case ecscomponents.EnemyBehaviorFleeing:
			stateFireRollPassed = nextRandomUnit(&enemyAIState.RNGState) < enemyAIFleeingFireChance
		}
	}

	switch enemyAIState.Behavior {
	case ecscomponents.EnemyBehaviorAggressive:
		intent.Fire = alignedToPlayer &&
			distanceToPlayer <= enemyAIAggressiveFireRange &&
			clearShot
	case ecscomponents.EnemyBehaviorCautious:
		intent.Fire = stateFireRollPassed &&
			alignedToPlayer &&
			distanceToPlayer <= enemyAICautiousFireRange &&
			clearShot
	case ecscomponents.EnemyBehaviorFleeing:
		intent.Fire = stateFireRollPassed &&
			alignedToPlayer &&
			distanceToPlayer <= enemyAIFleeingFireRange &&
			clearShot
	}

	if err := world.SetTankControlIntent(enemy, intent); err != nil {
		return fmt.Errorf("set enemy intent: %w", err)
	}

	if err := world.SetEnemyAI(enemy, enemyAIState); err != nil {
		return fmt.Errorf("set enemy ai state: %w", err)
	}

	return nil
}

func normalizedEnemyAIState(world *ecs.World, enemy ecs.Entity) ecscomponents.EnemyAI {
	enemyAIState, ok := world.EnemyAI(enemy)
	if !ok {
		return NewEnemyAI(enemy)
	}

	if enemyAIState.RNGState == 0 {
		enemyAIState.RNGState = seedForEnemyEntity(enemy)
	}

	if !isBehaviorValid(enemyAIState.FavoredBehavior) {
		enemyAIState.FavoredBehavior = sampleFavoredBehavior(&enemyAIState.RNGState)
	}

	if !isBehaviorValid(enemyAIState.Behavior) {
		enemyAIState.Behavior = enemyAIState.FavoredBehavior
	}

	if enemyAIState.PathIndex < 0 {
		enemyAIState.PathIndex = 0
	}
	if enemyAIState.PathIndex > len(enemyAIState.Path) {
		enemyAIState.PathIndex = len(enemyAIState.Path)
	}

	return enemyAIState
}

func isBehaviorValid(behavior ecscomponents.EnemyBehavior) bool {
	switch behavior {
	case ecscomponents.EnemyBehaviorAggressive, ecscomponents.EnemyBehaviorCautious, ecscomponents.EnemyBehaviorFleeing:
		return true
	default:
		return false
	}
}

func enemyHealthRatio(world *ecs.World, enemy ecs.Entity) float64 {
	health, ok := world.Health(enemy)
	if !ok || health.Max <= 0 {
		return 1
	}

	ratio := health.Current / health.Max
	return mathutil.Clamp(ratio, 0, 1)
}

func enemyCooldownReady(world *ecs.World, enemy ecs.Entity) bool {
	weapon, ok := world.Weapon(enemy)
	if !ok {
		return true
	}

	return weapon.CooldownRemaining <= 0
}

func behaviorThrottle(behavior ecscomponents.EnemyBehavior, clearShot bool, distanceToPlayer, stopDistance float64) float64 {
	switch behavior {
	case ecscomponents.EnemyBehaviorAggressive:
		if stopDistance > 0 && distanceToPlayer <= stopDistance {
			return 0
		}
		return 1
	case ecscomponents.EnemyBehaviorCautious:
		if clearShot {
			return 1
		}
		return 0.35
	case ecscomponents.EnemyBehaviorFleeing:
		return 1
	default:
		return 0
	}
}

func nextSteeringWaypoint(enemyAIState *ecscomponents.EnemyAI, transform ecscomponents.Transform) (x float64, y float64, ok bool) {
	for enemyAIState.PathIndex < len(enemyAIState.Path) {
		waypoint := enemyAIState.Path[enemyAIState.PathIndex]
		x = (float64(waypoint.Col) + 0.5) * enemyAIPathCellSize
		y = (float64(waypoint.Row) + 0.5) * enemyAIPathCellSize

		if distanceSquared(transform.X, transform.Y, x, y) <= enemyAIWaypointArrivalDistance*enemyAIWaypointArrivalDistance {
			enemyAIState.PathIndex++
			continue
		}

		return x, y, true
	}

	return 0, 0, false
}

func turnAndDistanceToPoint(transform ecscomponents.Transform, targetX, targetY float64) (turnDelta float64, distance float64) {
	dx := targetX - transform.X
	dy := targetY - transform.Y
	distanceSquaredValue := dx*dx + dy*dy
	if distanceSquaredValue <= 0.0001 {
		return 0, 0
	}

	desiredRotation := math.Atan2(dx, -dy)
	turnDelta = shortestAngleDelta(transform.RotationRadians, desiredRotation)

	return turnDelta, math.Sqrt(distanceSquaredValue)
}

func awayPoint(from ecscomponents.Transform, threat ecscomponents.Transform, distance float64) (x float64, y float64) {
	dx := from.X - threat.X
	dy := from.Y - threat.Y
	lengthSquared := dx*dx + dy*dy
	if lengthSquared <= 0.0001 {
		return from.X, from.Y + distance
	}

	invLength := 1 / math.Sqrt(lengthSquared)
	return from.X + dx*invLength*distance, from.Y + dy*invLength*distance
}

func markovTransitionWeights(current, favored ecscomponents.EnemyBehavior, healthRatio float64, cooldownReady bool, nearbyFriends int) behaviorWeights {
	weights := baseTransitionWeights(current)

	weights.applyFavored(favored)
	weights.applyHealth(healthRatio)
	weights.applyCooldown(cooldownReady)
	weights.applyFriends(nearbyFriends)

	weights.clampAndNormalize()
	return weights
}

func baseTransitionWeights(current ecscomponents.EnemyBehavior) behaviorWeights {
	switch current {
	case ecscomponents.EnemyBehaviorAggressive:
		return behaviorWeights{aggressive: 0.60, cautious: 0.30, fleeing: 0.10}
	case ecscomponents.EnemyBehaviorCautious:
		return behaviorWeights{aggressive: 0.35, cautious: 0.45, fleeing: 0.20}
	case ecscomponents.EnemyBehaviorFleeing:
		return behaviorWeights{aggressive: 0.20, cautious: 0.45, fleeing: 0.35}
	default:
		return behaviorWeights{aggressive: 0.60, cautious: 0.30, fleeing: 0.10}
	}
}

func (w *behaviorWeights) applyFavored(favored ecscomponents.EnemyBehavior) {
	switch favored {
	case ecscomponents.EnemyBehaviorAggressive:
		w.aggressive += 0.12
	case ecscomponents.EnemyBehaviorCautious:
		w.cautious += 0.12
	case ecscomponents.EnemyBehaviorFleeing:
		w.fleeing += 0.12
	}
}

func (w *behaviorWeights) applyHealth(healthRatio float64) {
	switch {
	case healthRatio > 0.66:
		w.aggressive += 0.20
		w.cautious += 0.05
		w.fleeing -= 0.15
	case healthRatio > 0.33:
		w.cautious += 0.10
	default:
		w.aggressive -= 0.20
		w.cautious += 0.05
		w.fleeing += 0.30
	}
}

func (w *behaviorWeights) applyCooldown(cooldownReady bool) {
	if cooldownReady {
		w.aggressive += 0.12
		w.cautious += 0.03
		w.fleeing -= 0.05
		return
	}

	w.aggressive -= 0.06
	w.cautious += 0.07
	w.fleeing += 0.04
}

func (w *behaviorWeights) applyFriends(nearbyFriends int) {
	if nearbyFriends <= 0 {
		return
	}

	friendCount := nearbyFriends
	if friendCount > enemyAIMaxFriendInfluenceCount {
		friendCount = enemyAIMaxFriendInfluenceCount
	}

	w.aggressive += 0.06 * float64(friendCount)
	w.cautious -= 0.01 * float64(friendCount)
	w.fleeing -= 0.03 * float64(friendCount)
}

func (w *behaviorWeights) clampAndNormalize() {
	w.aggressive = math.Max(w.aggressive, 0.01)
	w.cautious = math.Max(w.cautious, 0.01)
	w.fleeing = math.Max(w.fleeing, 0.01)

	total := w.aggressive + w.cautious + w.fleeing
	if total <= 0 {
		w.aggressive = 1
		w.cautious = 0
		w.fleeing = 0
		return
	}

	w.aggressive /= total
	w.cautious /= total
	w.fleeing /= total
}

func sampleBehavior(weights behaviorWeights, rngState *uint64) ecscomponents.EnemyBehavior {
	roll := nextRandomUnit(rngState)
	if roll < weights.aggressive {
		return ecscomponents.EnemyBehaviorAggressive
	}

	roll -= weights.aggressive
	if roll < weights.cautious {
		return ecscomponents.EnemyBehaviorCautious
	}

	return ecscomponents.EnemyBehaviorFleeing
}

func sampleFavoredBehavior(rngState *uint64) ecscomponents.EnemyBehavior {
	roll := nextRandomUnit(rngState)
	if roll < 0.5 {
		return ecscomponents.EnemyBehaviorAggressive
	}
	if roll < 0.8 {
		return ecscomponents.EnemyBehaviorCautious
	}

	return ecscomponents.EnemyBehaviorFleeing
}

func seedForEnemyEntity(entity ecs.Entity) uint64 {
	seed := uint64(entity)*0x9e3779b97f4a7c15 + 0xbf58476d1ce4e5b9
	if seed == 0 {
		seed = 0x94d049bb133111eb
	}

	return seed
}

func nextRandomUnit(rngState *uint64) float64 {
	n := nextRandomUint64(rngState)
	return float64(n>>11) / float64(uint64(1)<<53)
}

func nextRandomUint64(rngState *uint64) uint64 {
	if *rngState == 0 {
		*rngState = 0x94d049bb133111eb
	}

	x := *rngState
	x ^= x >> 12
	x ^= x << 25
	x ^= x >> 27
	*rngState = x

	return x * 2685821657736338717
}

func enemyColliderRadius(world *ecs.World, enemy ecs.Entity) float64 {
	collider, ok := world.CircleCollider(enemy)
	if !ok || collider.Radius <= 0 {
		return 16
	}

	return collider.Radius
}

func countNearbyEnemyFriends(world *ecs.World, enemy ecs.Entity, enemyTransform ecscomponents.Transform) int {
	count := 0
	radiusSquared := enemyAIFriendInfluenceRadius * enemyAIFriendInfluenceRadius

	for _, friend := range world.EnemyAITanks() {
		if friend.Entity == enemy {
			continue
		}

		if distanceSquared(enemyTransform.X, enemyTransform.Y, friend.Transform.X, friend.Transform.Y) <= radiusSquared {
			count++
			if count >= enemyAIMaxFriendInfluenceCount {
				return count
			}
		}
	}

	return count
}

func plannedTargetCellForBehavior(
	nav navGrid,
	behavior ecscomponents.EnemyBehavior,
	start navCell,
	targetTransform ecscomponents.Transform,
	obstacles []ecs.ColliderBody,
) navCell {
	playerCell := nav.worldToCell(targetTransform.X, targetTransform.Y)

	switch behavior {
	case ecscomponents.EnemyBehaviorCautious:
		if coverCell, ok := findNearestCoverCell(nav, start, targetTransform.X, targetTransform.Y, obstacles, enemyAICautiousCoverSearchRadius); ok {
			return coverCell
		}
		return start
	case ecscomponents.EnemyBehaviorFleeing:
		if fleeCell, ok := findBestFleeCell(nav, start, targetTransform.X, targetTransform.Y, obstacles, enemyAIFleeSearchRadius); ok {
			return fleeCell
		}
		return start
	default:
		if walkable, ok := nearestWalkableCell(nav, playerCell, max(nav.cols, nav.rows)); ok {
			return walkable
		}
		return start
	}
}

func buildNavGrid(bounds ecscomponents.MapBounds, obstacles []ecs.ColliderBody, agentRadius float64) navGrid {
	cols := int(math.Ceil(bounds.Width / enemyAIPathCellSize))
	rows := int(math.Ceil(bounds.Height / enemyAIPathCellSize))
	if cols < 1 {
		cols = 1
	}
	if rows < 1 {
		rows = 1
	}

	blocked := make([]bool, cols*rows)
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			cell := navCell{col: col, row: row}
			cx, cy := cellCenter(cell, enemyAIPathCellSize)
			for _, obstacle := range obstacles {
				r := obstacle.Collider.Radius + agentRadius + enemyAIPathPadding
				if distanceSquared(cx, cy, obstacle.Transform.X, obstacle.Transform.Y) <= r*r {
					blocked[row*cols+col] = true
					break
				}
			}
		}
	}

	return navGrid{
		cols:     cols,
		rows:     rows,
		cellSize: enemyAIPathCellSize,
		blocked:  blocked,
	}
}

func (g navGrid) worldToCell(x, y float64) navCell {
	col := int(math.Floor(x / g.cellSize))
	row := int(math.Floor(y / g.cellSize))

	if col < 0 {
		col = 0
	} else if col >= g.cols {
		col = g.cols - 1
	}

	if row < 0 {
		row = 0
	} else if row >= g.rows {
		row = g.rows - 1
	}

	return navCell{col: col, row: row}
}

func (g navGrid) inBounds(cell navCell) bool {
	return cell.col >= 0 && cell.col < g.cols && cell.row >= 0 && cell.row < g.rows
}

func (g navGrid) index(cell navCell) int {
	return cell.row*g.cols + cell.col
}

func (g navGrid) walkable(cell navCell) bool {
	if !g.inBounds(cell) {
		return false
	}

	return !g.blocked[g.index(cell)]
}

func findPathAStar(grid navGrid, start, goal navCell) []navCell {
	if !grid.inBounds(start) || !grid.inBounds(goal) {
		return nil
	}

	if !grid.walkable(goal) {
		var ok bool
		goal, ok = nearestWalkableCell(grid, goal, max(grid.cols, grid.rows))
		if !ok {
			return nil
		}
	}

	if start == goal {
		return []navCell{start}
	}

	totalCells := grid.cols * grid.rows
	unreachableCost := (int(^uint(0)>>1) / 4)
	gScore := make([]int, totalCells)
	cameFrom := make([]int, totalCells)
	closed := make([]bool, totalCells)
	inOpen := make([]bool, totalCells)

	for i := range gScore {
		gScore[i] = unreachableCost
		cameFrom[i] = -1
	}

	startIdx := grid.index(start)
	goalIdx := grid.index(goal)
	gScore[startIdx] = 0

	open := []int{startIdx}
	inOpen[startIdx] = true

	for len(open) > 0 {
		bestOpenIndex := 0
		bestIdx := open[0]
		bestCell := navCell{col: bestIdx % grid.cols, row: bestIdx / grid.cols}
		bestG := gScore[bestIdx]
		bestH := octileDistance(bestCell, goal)
		bestF := bestG + bestH

		for i := 1; i < len(open); i++ {
			idx := open[i]
			cell := navCell{col: idx % grid.cols, row: idx / grid.cols}
			g := gScore[idx]
			h := octileDistance(cell, goal)
			f := g + h

			if f < bestF ||
				(f == bestF && h < bestH) ||
				(f == bestF && h == bestH && cell.row < bestCell.row) ||
				(f == bestF && h == bestH && cell.row == bestCell.row && cell.col < bestCell.col) {
				bestOpenIndex = i
				bestIdx = idx
				bestCell = cell
				bestG = g
				bestH = h
				bestF = f
			}
		}

		open = append(open[:bestOpenIndex], open[bestOpenIndex+1:]...)
		inOpen[bestIdx] = false
		closed[bestIdx] = true

		if bestIdx == goalIdx {
			return reconstructPath(grid, cameFrom, goalIdx)
		}

		for _, step := range navNeighborSteps {
			neighbor := navCell{col: bestCell.col + step.dCol, row: bestCell.row + step.dRow}
			if !grid.inBounds(neighbor) {
				continue
			}

			if !grid.walkable(neighbor) && neighbor != goal {
				continue
			}

			if step.dCol != 0 && step.dRow != 0 {
				horizontal := navCell{col: bestCell.col + step.dCol, row: bestCell.row}
				vertical := navCell{col: bestCell.col, row: bestCell.row + step.dRow}
				if !grid.walkable(horizontal) || !grid.walkable(vertical) {
					continue
				}
			}

			neighborIdx := grid.index(neighbor)
			if closed[neighborIdx] {
				continue
			}

			tentativeG := bestG + step.cost
			if tentativeG >= gScore[neighborIdx] {
				continue
			}

			cameFrom[neighborIdx] = bestIdx
			gScore[neighborIdx] = tentativeG

			if !inOpen[neighborIdx] {
				open = append(open, neighborIdx)
				inOpen[neighborIdx] = true
			}
		}
	}

	return nil
}

func reconstructPath(grid navGrid, cameFrom []int, goalIdx int) []navCell {
	pathReversed := make([]navCell, 0, 8)
	for idx := goalIdx; idx >= 0; idx = cameFrom[idx] {
		pathReversed = append(pathReversed, navCell{col: idx % grid.cols, row: idx / grid.cols})
		if cameFrom[idx] < 0 {
			break
		}
	}

	path := make([]navCell, 0, len(pathReversed))
	for i := len(pathReversed) - 1; i >= 0; i-- {
		path = append(path, pathReversed[i])
	}

	return path
}

func octileDistance(a, b navCell) int {
	dCol := absInt(a.col - b.col)
	dRow := absInt(a.row - b.row)
	if dCol > dRow {
		return 14*dRow + 10*(dCol-dRow)
	}

	return 14*dCol + 10*(dRow-dCol)
}

func nearestWalkableCell(grid navGrid, origin navCell, maxRadius int) (navCell, bool) {
	if grid.walkable(origin) {
		return origin, true
	}

	if maxRadius < 1 {
		maxRadius = 1
	}

	for radius := 1; radius <= maxRadius; radius++ {
		best := navCell{}
		hasBest := false
		bestDistSq := int(^uint(0) >> 1)

		for row := origin.row - radius; row <= origin.row+radius; row++ {
			for col := origin.col - radius; col <= origin.col+radius; col++ {
				cell := navCell{col: col, row: row}
				if !grid.inBounds(cell) || !grid.walkable(cell) {
					continue
				}

				dCol := col - origin.col
				dRow := row - origin.row
				distSq := dCol*dCol + dRow*dRow
				if !hasBest || distSq < bestDistSq || (distSq == bestDistSq && (row < best.row || (row == best.row && col < best.col))) {
					best = cell
					hasBest = true
					bestDistSq = distSq
				}
			}
		}

		if hasBest {
			return best, true
		}
	}

	return navCell{}, false
}

func findNearestCoverCell(
	grid navGrid,
	start navCell,
	targetX, targetY float64,
	obstacles []ecs.ColliderBody,
	searchRadius int,
) (navCell, bool) {
	best := navCell{}
	hasBest := false
	bestDistSq := int(^uint(0) >> 1)

	for row := start.row - searchRadius; row <= start.row+searchRadius; row++ {
		for col := start.col - searchRadius; col <= start.col+searchRadius; col++ {
			cell := navCell{col: col, row: row}
			if !grid.inBounds(cell) || !grid.walkable(cell) {
				continue
			}

			x, y := cellCenter(cell, grid.cellSize)
			if lineOfSightClear(x, y, targetX, targetY, obstacles) {
				continue
			}

			dCol := col - start.col
			dRow := row - start.row
			distSq := dCol*dCol + dRow*dRow
			if !hasBest || distSq < bestDistSq || (distSq == bestDistSq && (row < best.row || (row == best.row && col < best.col))) {
				best = cell
				hasBest = true
				bestDistSq = distSq
			}
		}
	}

	return best, hasBest
}

func findBestFleeCell(
	grid navGrid,
	start navCell,
	targetX, targetY float64,
	obstacles []ecs.ColliderBody,
	searchRadius int,
) (navCell, bool) {
	best := navCell{}
	hasBest := false
	bestScore := -1.0

	for row := start.row - searchRadius; row <= start.row+searchRadius; row++ {
		for col := start.col - searchRadius; col <= start.col+searchRadius; col++ {
			cell := navCell{col: col, row: row}
			if !grid.inBounds(cell) || !grid.walkable(cell) {
				continue
			}

			x, y := cellCenter(cell, grid.cellSize)
			score := distanceSquared(x, y, targetX, targetY)
			if !lineOfSightClear(x, y, targetX, targetY, obstacles) {
				score += enemyAIFleeCoverBonus
			}

			if !hasBest || score > bestScore || (score == bestScore && (row < best.row || (row == best.row && col < best.col))) {
				best = cell
				hasBest = true
				bestScore = score
			}
		}
	}

	if hasBest {
		return best, true
	}

	return start, false
}

func lineOfSightClear(ax, ay, bx, by float64, obstacles []ecs.ColliderBody) bool {
	segmentDx := bx - ax
	segmentDy := by - ay
	segmentLengthSquared := segmentDx*segmentDx + segmentDy*segmentDy

	minX := math.Min(ax, bx)
	maxX := math.Max(ax, bx)
	minY := math.Min(ay, by)
	maxY := math.Max(ay, by)

	for _, obstacle := range obstacles {
		r := obstacle.Collider.Radius
		ox := obstacle.Transform.X
		oy := obstacle.Transform.Y

		if ox < minX-r || ox > maxX+r || oy < minY-r || oy > maxY+r {
			continue
		}

		if segmentLengthSquared <= 0.000001 {
			if distanceSquared(ax, ay, ox, oy) <= r*r {
				return false
			}
			continue
		}

		t := ((ox-ax)*segmentDx + (oy-ay)*segmentDy) / segmentLengthSquared
		t = mathutil.Clamp(t, 0, 1)

		closestX := ax + t*segmentDx
		closestY := ay + t*segmentDy
		if distanceSquared(closestX, closestY, ox, oy) <= r*r {
			return false
		}
	}

	return true
}

func distanceSquared(ax, ay, bx, by float64) float64 {
	dx := ax - bx
	dy := ay - by
	return dx*dx + dy*dy
}

func cellCenter(cell navCell, cellSize float64) (x float64, y float64) {
	return (float64(cell.col) + 0.5) * cellSize, (float64(cell.row) + 0.5) * cellSize
}

func absInt(v int) int {
	if v < 0 {
		return -v
	}

	return v
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func shortestAngleDelta(current, target float64) float64 {
	delta := math.Mod((target-current)+math.Pi, 2*math.Pi)
	if delta < 0 {
		delta += 2 * math.Pi
	}

	return delta - math.Pi
}
