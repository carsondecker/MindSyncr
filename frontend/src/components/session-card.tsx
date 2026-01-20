import { Button } from "@/components/ui/button"
import { Card, CardDescription, CardFooter, CardHeader, CardTitle } from "./ui/card"
import { Badge } from "./ui/badge"
import type { Session } from "@/lib/api/models/sessions"
import { LogIn, LogOut, Play, StopCircle, Trash2 } from "lucide-react"
import { useNavigate } from "react-router-dom"

type SessionCardProps = {
    session: Session
    onDelete?: (session_id: string) => void
    onEnd?: (session_id: string) => void
    onJoin?: (session_id: string) => void
    onLeave?: (session_id: string) => void
}

export function SessionCard({ session, onDelete, onEnd, onJoin, onLeave } : SessionCardProps) {
  const navigate = useNavigate();
  
  const onOpenOwner = (id: string) => {
    navigate(`/sessions/${id}/presenter`)
  }

  const onOpenMember = (id: string) => {
    navigate(`/sessions/${id}/viewer`)
  }

  return (
    <Card className="hover:shadow-lg transition-shadow border-l-4 border-l-green-500">
      <CardHeader>
        <div className="flex items-start justify-between">
          <div className="flex-1">
            <CardTitle className="text-lg">{session.name}</CardTitle>
          </div>
          <div className="flex flex-col gap-2 items-end">
            {session.is_active && (
              <Badge className="bg-green-500">Active</Badge>
            )}
            {session.is_member && (
              <Badge variant="outline">Joined</Badge>
            )}
          </div>
        </div>
        <p className="text-xs text-gray-500 mt-2">Created: {new Date(session.created_at).toLocaleDateString()}</p>
      </CardHeader>
      <CardFooter className="flex gap-2 justify-end flex-wrap">
        {session.is_owner ? (
          <>
            <Button
              variant="outline"
              size="sm"
              onClick={() => onOpenOwner(session.id)}
              className="gap-2"
            >
              <Play className="h-4 w-4" />
              Open
            </Button>
            {session.is_active && (
              <Button
                variant="outline"
                size="sm"
                onClick={() => onEnd?.(session.id)}
                className="gap-2 hover:bg-orange-50 hover:text-orange-600 hover:border-orange-300"
              >
                <StopCircle className="h-4 w-4" />
                End Session
              </Button>
            )}
            <Button
              variant="destructive"
              size="sm"
              onClick={() => onDelete?.(session.id)}
              className="gap-2"
            >
              <Trash2 className="h-4 w-4" />
              Delete
            </Button>
          </>
        ) : (
          <>
            {session.is_active && !session.is_member && (
              <Button
                variant="default"
                size="sm"
                onClick={() => onJoin?.(session.id)}
                className="gap-2"
              >
                <LogIn className="h-4 w-4" />
                Join Session
              </Button>
            )}
            {session.is_member && (
              <>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => onOpenMember(session.id)}
                  className="gap-2"
                >
                  <Play className="h-4 w-4" />
                  Open
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => onLeave?.(session.id)}
                  className="gap-2 hover:bg-red-50 hover:text-red-600 hover:border-red-300"
                >
                  <LogOut className="h-4 w-4" />
                  Leave
                </Button>
              </>
              
            )}
          </>
        )}
      </CardFooter>
    </Card>
  )
}
