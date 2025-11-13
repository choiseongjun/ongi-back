package utils

import (
	"math"
	"sync"
)

// Vector5D - 5차원 벡터 (성향 점수)
type Vector5D struct {
	Sociality   float64
	Activity    float64
	Intimacy    float64
	Immersion   float64
	Flexibility float64
}

// ToSlice - 벡터를 슬라이스로 변환
func (v *Vector5D) ToSlice() []float64 {
	return []float64{v.Sociality, v.Activity, v.Intimacy, v.Immersion, v.Flexibility}
}

// FromSlice - 슬라이스에서 벡터 생성
func FromSlice(s []float64) *Vector5D {
	if len(s) != 5 {
		return nil
	}
	return &Vector5D{
		Sociality:   s[0],
		Activity:    s[1],
		Intimacy:    s[2],
		Immersion:   s[3],
		Flexibility: s[4],
	}
}

// Magnitude - 벡터의 크기 계산
func (v *Vector5D) Magnitude() float64 {
	return math.Sqrt(
		v.Sociality*v.Sociality +
			v.Activity*v.Activity +
			v.Intimacy*v.Intimacy +
			v.Immersion*v.Immersion +
			v.Flexibility*v.Flexibility,
	)
}

// EuclideanDistance - 유클리드 거리 계산
func EuclideanDistance(v1, v2 *Vector5D) float64 {
	diff1 := v1.Sociality - v2.Sociality
	diff2 := v1.Activity - v2.Activity
	diff3 := v1.Intimacy - v2.Intimacy
	diff4 := v1.Immersion - v2.Immersion
	diff5 := v1.Flexibility - v2.Flexibility

	return math.Sqrt(
		diff1*diff1 + diff2*diff2 + diff3*diff3 + diff4*diff4 + diff5*diff5,
	)
}

// CosineSimilarity - 코사인 유사도 계산 (더 정확함)
func CosineSimilarity(v1, v2 *Vector5D) float64 {
	dotProduct := v1.Sociality*v2.Sociality +
		v1.Activity*v2.Activity +
		v1.Intimacy*v2.Intimacy +
		v1.Immersion*v2.Immersion +
		v1.Flexibility*v2.Flexibility

	magnitude1 := v1.Magnitude()
	magnitude2 := v2.Magnitude()

	if magnitude1 == 0 || magnitude2 == 0 {
		return 0
	}

	return dotProduct / (magnitude1 * magnitude2)
}

// Similarity - 유사도를 0-100 점수로 변환
func Similarity(v1, v2 *Vector5D) float64 {
	distance := EuclideanDistance(v1, v2)
	maxDistance := math.Sqrt(5 * 100 * 100) // 최대 거리 (각 차원 최대 100)
	similarity := (1 - (distance / maxDistance)) * 100

	// 음수 방지
	if similarity < 0 {
		similarity = 0
	}

	return math.Round(similarity*10) / 10
}

// SimilarityScore - 코사인 유사도를 0-100 점수로 변환
func SimilarityScore(v1, v2 *Vector5D) float64 {
	cosineSim := CosineSimilarity(v1, v2)
	// 코사인 유사도는 -1 ~ 1 범위, 0 ~ 1로 정규화 후 100배
	score := ((cosineSim + 1) / 2) * 100
	return math.Round(score*10) / 10
}

// BatchSimilarity - 여러 벡터와의 유사도를 병렬로 계산
type SimilarityResult struct {
	Index      int
	Similarity float64
}

func BatchSimilarity(target *Vector5D, vectors []*Vector5D, workers int) []SimilarityResult {
	if workers <= 0 {
		workers = 4 // 기본 워커 수
	}

	results := make([]SimilarityResult, len(vectors))
	var wg sync.WaitGroup

	// 작업을 워커들에게 분배
	chunkSize := len(vectors) / workers
	if chunkSize == 0 {
		chunkSize = 1
	}

	for w := 0; w < workers; w++ {
		start := w * chunkSize
		end := start + chunkSize

		if w == workers-1 {
			end = len(vectors) // 마지막 워커는 나머지 전부 처리
		}

		if start >= len(vectors) {
			break
		}

		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i < end && i < len(vectors); i++ {
				results[i] = SimilarityResult{
					Index:      i,
					Similarity: SimilarityScore(target, vectors[i]),
				}
			}
		}(start, end)
	}

	wg.Wait()
	return results
}

// WeightedVector - 가중치를 적용한 벡터
func (v *Vector5D) ApplyWeights(weights *Vector5D) *Vector5D {
	return &Vector5D{
		Sociality:   v.Sociality * weights.Sociality,
		Activity:    v.Activity * weights.Activity,
		Intimacy:    v.Intimacy * weights.Intimacy,
		Immersion:   v.Immersion * weights.Immersion,
		Flexibility: v.Flexibility * weights.Flexibility,
	}
}

// Normalize - 벡터 정규화 (단위 벡터로 변환)
func (v *Vector5D) Normalize() *Vector5D {
	mag := v.Magnitude()
	if mag == 0 {
		return &Vector5D{}
	}

	return &Vector5D{
		Sociality:   v.Sociality / mag,
		Activity:    v.Activity / mag,
		Intimacy:    v.Intimacy / mag,
		Immersion:   v.Immersion / mag,
		Flexibility: v.Flexibility / mag,
	}
}

// ManhattanDistance - 맨해튼 거리 (택시 거리)
func ManhattanDistance(v1, v2 *Vector5D) float64 {
	return math.Abs(v1.Sociality-v2.Sociality) +
		math.Abs(v1.Activity-v2.Activity) +
		math.Abs(v1.Intimacy-v2.Intimacy) +
		math.Abs(v1.Immersion-v2.Immersion) +
		math.Abs(v1.Flexibility-v2.Flexibility)
}
