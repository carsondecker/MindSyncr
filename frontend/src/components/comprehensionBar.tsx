import type { ComprehensionScore } from "@/lib/api/models/comprehensionScores"
import { Card, CardContent, CardHeader, CardTitle } from "./ui/card"
import { useMemo } from "react"

type ComprehensionBarProps = {
  userScores: Record<string, ComprehensionScore>
  showCounts?: boolean
}

export default function ComprehensionBar({ userScores, showCounts }: ComprehensionBarProps) {
  const { scoreCounts, totalUsers, scorePercentages } = useMemo(() => {
    console.log(userScores)

    const counts: Record<number, number> = { 1: 0, 2: 0, 3: 0, 4: 0, 5: 0 }
    
    Object.values(userScores || {}).forEach(scoreObj => {
      if (scoreObj.score >= 1 && scoreObj.score <= 5) {
        counts[scoreObj.score]++
      }
    })

    const total = Object.keys(userScores || {}).length
    
    const percentages: Record<number, number> = {
      1: total > 0 ? (counts[1] / total) * 100 : 0,
      2: total > 0 ? (counts[2] / total) * 100 : 0,
      3: total > 0 ? (counts[3] / total) * 100 : 0,
      4: total > 0 ? (counts[4] / total) * 100 : 0,
      5: total > 0 ? (counts[5] / total) * 100 : 0,
    }

    return { scoreCounts: counts, totalUsers: total, scorePercentages: percentages }
  }, [userScores])

  const scoreColors = {
    1: "bg-red-500",
    2: "bg-orange-500",
    3: "bg-yellow-500",
    4: "bg-lime-500",
    5: "bg-green-500",
  }

  const scoreLabels = {
    1: "Very Low",
    2: "Low",
    3: "Medium",
    4: "Good",
    5: "Excellent",
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center justify-between">
          <span>Comprehension Scores</span>
          <span className="text-sm font-normal text-gray-500">
            {totalUsers} {totalUsers === 1 ? 'student' : 'students'}
          </span>
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          <div className="flex h-12 w-full rounded-lg overflow-hidden border border-gray-200">
            {([1, 2, 3, 4, 5] as const).map((score) => {
              const percentage = scorePercentages[score]
              if (percentage === 0) return null
              
              return (
                <div
                  key={score}
                  className={`${scoreColors[score]} flex items-center justify-center text-white font-semibold transition-all hover:opacity-80 relative group`}
                  style={{ width: `${percentage}%` }}
                >
                  {percentage >= 10 && (
                    <span className="text-sm">
                      {Math.round(percentage)}%
                    </span>
                  )}
                  
                  <div className="absolute bottom-full mb-2 hidden group-hover:block bg-gray-900 text-white text-xs rounded py-1 px-2 whitespace-nowrap z-10">
                    Score {score}: {scoreCounts[score]} {scoreCounts[score] === 1 ? 'student' : 'students'}
                  </div>
                </div>
              )
            })}
          </div>

          {showCounts && (
            <div className="grid grid-cols-5 gap-2 text-sm">
              {([1, 2, 3, 4, 5] as const).map((score) => (
                <div key={score} className="flex flex-col items-center">
                  <div className={`w-4 h-4 rounded ${scoreColors[score]} mb-1`} />
                  <span className="text-xs font-medium">{score}</span>
                  <span className="text-xs text-gray-500">{scoreLabels[score]}</span>
                  <span className="text-xs font-semibold text-gray-700">
                    {scoreCounts[score]}
                  </span>
                </div>
              ))}
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  )
}