import { useState, useMemo } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "./ui/card"
import { Button } from "./ui/button"
import { Textarea } from "./ui/textarea"
import { ThumbsUp, MessageCircle, Flag, Trash2, Edit2, Check } from "lucide-react"
import type { Question } from "@/lib/api/models/questions"

type QuestionsProps = {
  questions: Record<string, Question>,
  userId: string,
  onQuestionSubmit?: (text: string) => void
  onQuestionDelete?: (id: string) => void
}

type SortOption = "latest" | "popular"

// TODO: Add version that presenters cannot send questions, may need to move adding questions to a new component
export default function Questions({ questions, userId, onQuestionSubmit, onQuestionDelete }: QuestionsProps) {
  const [newQuestion, setNewQuestion] = useState("")
  const [sortBy, setSortBy] = useState<SortOption>("latest")
  const [editingQuestion, setEditingQuestion] = useState<string | null>(null)
  const [editText, setEditText] = useState("")

  const sortedQuestions = useMemo(() => {
    const sorted = Object.keys(questions).map(key => questions[key])
    
    if (sortBy === "latest") {
      return sorted.sort((a, b) => b.created_at.getTime() - a.created_at.getTime())
    } else {
      // TODO: Implement popular sorting when likes are added
      return sorted.sort((a, b) => b.created_at.getTime() - a.created_at.getTime())
    }
  }, [questions, sortBy])

  const handleSubmitQuestion = (text: string) => {
  if (text && onQuestionSubmit) {
    onQuestionSubmit(text)
    setNewQuestion("") // Clear the input
  }
}

  const handleLikeQuestion = (id: string) => {
    // TODO: Implement like functionality when added
  }

  const handleMarkAnswered = (id: string) => {
    // TODO: Implement mark as answered
  }

  const handleReplyToQuestion = (id: string) => {
    // TODO: Implement reply functionality when added
  }

  const handleReportQuestion = (id: string) => {
    // TODO: Implement report functionality when added
  }

  const handleDeleteQuestion = (id: string) => {
    if (onQuestionDelete) {
      onQuestionDelete(id)
  }
  }

  const handleStartEdit = (question: Question) => {
    setEditingQuestion(question.id)
    setEditText(question.text)
  }

  const handleSaveEdit = (questionId: string) => {
    // TODO: Implement edit question
    setEditingQuestion(null)
    setEditText("")
  }

  const handleCancelEdit = () => {
    setEditingQuestion(null)
    setEditText("")
  }

  return (
    <Card className="mt-6 w-full">
      <CardHeader>
        <CardTitle className="flex items-center justify-between">
          <span>Questions & Answers</span>
          <div className="flex gap-2">
            <Button
              variant={sortBy === "latest" ? "default" : "outline"}
              size="sm"
              onClick={() => setSortBy("latest")}
            >
              Latest
            </Button>
            <Button
              variant={sortBy === "popular" ? "default" : "outline"}
              size="sm"
              onClick={() => setSortBy("popular")}
              disabled
              title="Coming soon"
            >
              Most Popular
            </Button>
          </div>
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* New Question Input */}
        <div className="space-y-2">
          <Textarea
            placeholder="Ask a question..."
            value={newQuestion}
            onChange={(e) => setNewQuestion(e.target.value)}
            className="min-h-[80px]"
          />
          <div className="flex justify-end">
            <Button onClick={() => handleSubmitQuestion(newQuestion)} disabled={!newQuestion.trim()}>
              Post Question
            </Button>
          </div>
        </div>

        {/* Questions List */}
        <div className="space-y-4">
          {sortedQuestions.length === 0 ? (
            <p className="text-center text-gray-500 py-8">
              No questions yet. Be the first to ask!
            </p>
          ) : (
            sortedQuestions.map((question) => {
              const isOwner = question.user_id === userId
              const isEditing = editingQuestion === question.id

              return (
                <div
                  key={question.id}
                  className={`border rounded-lg p-4 space-y-3 ${
                    question.is_answered ? "bg-green-50 border-green-200" : ""
                  }`}
                >
                  {/* Question Header */}
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <div className="flex items-center gap-2 mb-1">
                        <span className="text-xs text-gray-500">
                          {question.created_at.toLocaleString()}
                        </span>
                        {question.is_answered && (
                          <span className="flex items-center gap-1 text-xs bg-green-500 text-white px-2 py-0.5 rounded-full">
                            <Check className="w-3 h-3" />
                            Answered
                          </span>
                        )}
                      </div>
                      
                      {/* Question Text */}
                      {isEditing ? (
                        <div className="space-y-2">
                          <Textarea
                            value={editText}
                            onChange={(e) => setEditText(e.target.value)}
                            className="min-h-[60px]"
                          />
                          <div className="flex gap-2">
                            <Button
                              size="sm"
                              onClick={() => handleSaveEdit(question.id)}
                              disabled={!editText.trim()}
                            >
                              Save
                            </Button>
                            <Button
                              size="sm"
                              variant="outline"
                              onClick={handleCancelEdit}
                            >
                              Cancel
                            </Button>
                          </div>
                        </div>
                      ) : (
                        <p className="text-gray-800">{question.text}</p>
                      )}
                    </div>
                  </div>

                  {/* Question Actions */}
                  <div className="flex items-center gap-2 text-sm">
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => handleLikeQuestion(question.id)}
                      disabled
                      title="Coming soon"
                    >
                      <ThumbsUp className="w-4 h-4 mr-1" />
                      0
                    </Button>

                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => handleReplyToQuestion(question.id)}
                      disabled
                      title="Coming soon"
                    >
                      <MessageCircle className="w-4 h-4 mr-1" />
                      0
                    </Button>

                    {isOwner && (
                      <>
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => handleMarkAnswered(question.id)}
                        >
                          <Check className="w-4 h-4 mr-1" />
                          {question.is_answered ? "Unmark" : "Mark Answered"}
                        </Button>

                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => handleStartEdit(question)}
                        >
                          <Edit2 className="w-4 h-4 mr-1" />
                          Edit
                        </Button>

                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => handleDeleteQuestion(question.id)}
                          className="text-red-600 hover:text-red-700"
                        >
                          <Trash2 className="w-4 h-4 mr-1" />
                          Delete
                        </Button>
                      </>
                    )}

                    {!isOwner && (
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => handleReportQuestion(question.id)}
                        className="text-orange-600 hover:text-orange-700"
                        disabled
                        title="Coming soon"
                      >
                        <Flag className="w-4 h-4 mr-1" />
                        Report
                      </Button>
                    )}
                  </div>
                </div>
              )
            })
          )}
        </div>
      </CardContent>
    </Card>
  )
}