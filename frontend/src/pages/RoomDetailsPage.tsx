import { SessionCard } from "@/components/session-card"
import { Button } from "@/components/ui/button"
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import { useAuth } from "@/lib/context/AuthContext"
import useRooms from "@/lib/hooks/useRooms"
import useSessionMutations from "@/lib/hooks/useSessionMutations"
import useSessions from "@/lib/hooks/useSessions"
import { ArrowLeft, Plus } from "lucide-react"
import { useEffect, useState } from "react"
import { useNavigate, useParams } from "react-router-dom"


export default function RoomDetailsPage() {
    const { room_id } = useParams()
    
    const { fetchRoomById, room } = useRooms()
    const { sessions } = useSessions(room_id!)
    const { createSession, deleteSession, endSession, joinSession, leaveSession } = useSessionMutations(room_id!)

    const navigate = useNavigate()
    const [showCreateSession, setShowCreateSession] = useState(false)
    const [sessionName, setSessionName] = useState("")
    const [sessionDescription, setSessionDescription] = useState("")

    const { user } = useAuth()

    useEffect(() => {
        fetchRoomById(room_id!)
    }, [])

    const handleCreateSession = () => {
        setShowCreateSession(true)
    }

    const handleCreateSessionSubmit = () => {
        setShowCreateSession(false)
        createSession.mutate({ name: sessionName })
    }

    const handleDeleteSession = (session_id: string) => {
        deleteSession.mutate(session_id)
    }

    const handleEndSession = (session_id: string) => {
        endSession.mutate(session_id)
    }

    const handleJoinSession = (session_id: string) => {
        joinSession.mutate(session_id)
    }

    const handleLeaveSession = (session_id: string) => {
        leaveSession.mutate(session_id)
    }

    const handleBackToRooms = () => {
        navigate("/")
    }

    if (room.isPending || sessions.isPending) {
        return <div>Loading…</div>
    }

    if (room.isError) {
        return <div>Error: {room.error.message}</div>
    }
    
    if (sessions.isError) {
        return <div>Error: {sessions.error.message}</div>
    }

    const isOwner = room.data.owner_id === user?.id
    const hasActiveSession = sessions.data.some(s => s.is_active)

    return (
        <div className="max-w-5xl mx-auto p-6">
            <Button
                variant="ghost"
                onClick={handleBackToRooms}
                className="mb-4 gap-2"
            >
                <ArrowLeft className="h-4 w-4" />
                Back to Rooms
            </Button>

            <div className="mb-8">
                <div className="flex items-start justify-between mb-4">
                <div>
                    <h1 className="text-3xl font-bold text-gray-900 mb-2">{room.data.name}</h1>
                    {room.data.description && (
                    <p className="text-gray-600 mb-2">{room.data.description}</p>
                    )}
                    <div className="flex gap-4 text-sm text-gray-500">
                    <span>Room ID: {room.data.id}</span>
                    <span>Join Code: <span className="font-mono font-semibold text-blue-600">{room.data.join_code}</span></span>
                    </div>
                </div>
                </div>
            </div>

            <div className="mb-6 flex items-center justify-between">
                <h2 className="text-2xl font-semibold text-gray-900">Sessions</h2>
                {isOwner && (
                <Button
                    onClick={handleCreateSession}
                    className="gap-2"
                    disabled={hasActiveSession}
                >
                    <Plus className="h-4 w-4" />
                    Create Session
                </Button>
                )}
            </div>

            {hasActiveSession && isOwner && (
                <div className="mb-4 p-4 bg-amber-50 border border-amber-200 rounded-lg">
                <p className="text-sm text-amber-800">
                    <strong>Note:</strong> You have an active session. End it before creating a new one.
                </p>
                </div>
            )}

            <div className="space-y-3">
                {sessions.data && sessions.data.length > 0 ? (
                sessions.data.map((session) => (
                    <SessionCard
                    key={session.id}
                    session={session}
                    onDelete={handleDeleteSession}
                    onEnd={handleEndSession}
                    onJoin={handleJoinSession}
                    onLeave={handleLeaveSession}
                    />
                ))
                ) : (
                <div className="text-center py-16 bg-gray-50 rounded-lg border-2 border-dashed border-gray-300">
                    <p className="text-gray-500 text-lg">No sessions yet</p>
                    <p className="text-gray-400 text-sm mt-2">
                    {isOwner
                        ? "Create your first session to get started"
                        : "Wait for the instructor to create a session"}
                    </p>
                </div>
                )}
            </div>

            <Dialog open={showCreateSession} onOpenChange={setShowCreateSession}>
                <DialogContent className="sm:max-w-md">
                <DialogHeader>
                    <DialogTitle>Create New Session</DialogTitle>
                    <DialogDescription>
                    Add a new learning session for this room
                    </DialogDescription>
                </DialogHeader>
                <div className="space-y-4 py-4">
                    <div className="space-y-2">
                    <Label htmlFor="session-name">Session Name *</Label>
                    <Input
                        id="session-name"
                        placeholder="e.g., Week 1: Introduction"
                        value={sessionName}
                        onChange={(e) => setSessionName(e.target.value)}
                    />
                    </div>
                    <div className="space-y-2">
                    <Label htmlFor="session-description">Description (optional)</Label>
                    <Textarea
                        id="session-description"
                        placeholder="Add details about this session..."
                        value={sessionDescription}
                        onChange={(e) => setSessionDescription(e.target.value)}
                        rows={3}
                    />
                    </div>
                </div>
                <DialogFooter>
                    <Button
                    variant="outline"
                    onClick={() => setShowCreateSession(false)}
                    >
                    Cancel
                    </Button>
                    <Button
                    onClick={handleCreateSessionSubmit}
                    disabled={!sessionName.trim()}
                    >
                    Create Session
                    </Button>
                </DialogFooter>
                </DialogContent>
            </Dialog>
        </div>
  )
}