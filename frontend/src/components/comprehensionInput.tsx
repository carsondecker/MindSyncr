import { useState } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"

type ComprehensionInputProps = {
  currentScore?: number
  onScoreSubmit: (score: number) => void
  disabled?: boolean
}

export default function ComprehensionInput({ 
  currentScore,
  onScoreSubmit,
  disabled = false
}: ComprehensionInputProps) {
  const [selectedScore, setSelectedScore] = useState<number | null>(currentScore || null)

  const handleScoreClick = (score: number) => {
    if (selectedScore !== score) {
      setSelectedScore(score)
      onScoreSubmit(score)
    }
  }

  const scoreColors = {
    1: { bg: "bg-red-500", hover: "hover:bg-red-600", border: "border-red-500", text: "text-red-500" },
    2: { bg: "bg-orange-500", hover: "hover:bg-orange-600", border: "border-orange-500", text: "text-orange-500" },
    3: { bg: "bg-yellow-500", hover: "hover:bg-yellow-600", border: "border-yellow-500", text: "text-yellow-500" },
    4: { bg: "bg-lime-500", hover: "hover:bg-lime-600", border: "border-lime-500", text: "text-lime-500" },
    5: { bg: "bg-green-500", hover: "hover:bg-green-600", border: "border-green-500", text: "text-green-500" },
  }

  const scoreLabels = {
    1: "Very Low",
    2: "Low",
    3: "Medium",
    4: "Good",
    5: "Excellent",
  }

  const scoreEmojis = {
    1: "😟",
    2: "😕",
    3: "😐",
    4: "🙂",
    5: "😄",
  }

  return (
    <Card className="w-full max-w-2xl">
      <CardHeader className="pb-3">
        <CardTitle className="text-lg">How well do you understand?</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-3">
          <div className="grid grid-cols-5 gap-2">
            {([1, 2, 3, 4, 5] as const).map((score) => {
              const isSelected = selectedScore === score
              const colors = scoreColors[score]
              
              return (
                <button
                  key={score}
                  onClick={() => handleScoreClick(score)}
                  disabled={disabled}
                  className={`
                    flex flex-col items-center justify-center p-3 rounded-lg border-2 
                    transition-all duration-200 
                    ${isSelected 
                      ? `${colors.bg} ${colors.border} text-white shadow-lg scale-105` 
                      : `bg-white ${colors.border} ${colors.text} ${colors.hover} hover:scale-105 hover:shadow-md`
                    }
                    ${disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'}
                  `}
                >
                  <span className="text-2xl mb-1">{scoreEmojis[score]}</span>
                  <span className="text-xl font-bold">{score}</span>
                  <span className={`text-[10px] font-medium ${isSelected ? 'text-white' : 'text-gray-600'}`}>
                    {scoreLabels[score]}
                  </span>
                </button>
              )
            })}
          </div>
          {selectedScore && (
            <div className="text-center p-2 bg-blue-50 border border-blue-200 rounded-lg">
              <p className="text-xs text-blue-800">
                <span className="font-semibold">Your level:</span> {scoreLabels[selectedScore as keyof typeof scoreLabels]} ({selectedScore}/5)
              </p>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  )
}