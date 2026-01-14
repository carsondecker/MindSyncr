import { RoomCard } from "@/components/room-card"
import { Button } from "@/components/ui/button"
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { type CreateRoomRequest, type Room } from "@/lib/api/models/rooms"
import { useAuth } from "@/lib/context/AuthContext"
import useRoomsApi from "@/lib/hooks/useRoomsApi"
import { useMutation, useQuery } from "@tanstack/react-query"
import { useState } from "react"
import { Plus, UserPlus, Camera } from "lucide-react"
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog"
import { Label } from "@/components/ui/label"
import { Input } from "@/components/ui/input"
import { Textarea } from "@/components/ui/textarea"

export default function HomePage() {
    const { user } = useAuth()
    const { getJoinedRooms, getOwnedRooms, createRoom, deleteRoom, joinRoom, leaveRoom } = useRoomsApi()
    const [activeTab, setActiveTab] = useState<"owned" | "joined">("owned")
    const [showCreateRoom, setShowCreateRoom] = useState(false)
    const [showJoinRoom, setShowJoinRoom] = useState(false)
    const [joinMethod, setJoinMethod] = useState<"code" | "qr">("code")
    const [joinCode, setJoinCode] = useState("")
    const [roomName, setRoomName] = useState("")
    const [roomDescription, setRoomDescription] = useState("")


    const ownedRoomsQuery = useQuery({
        queryKey: ['rooms', 'owned'],
        queryFn: getOwnedRooms
    })

    const joinedRoomsQuery = useQuery({
        queryKey: ['rooms', 'joined'],
        queryFn: getJoinedRooms
    })

    const addRoomQuery = useMutation({
        mutationKey: ['addRoom'],
        mutationFn: (data: CreateRoomRequest) => createRoom(data),
        /*
        onMutate: async (newRoom, context) => {
                await context.client.cancelQueries({ queryKey: ['rooms'] })

            const prevRooms = context.client.getQueryData(['rooms'])

            context.client.setQueryData(['rooms'], (old: Room[]) => [...old, newRoom])

        return { prevRooms }
        },
        */
        onError: (err, variables, onMutateResult, context) => {
            console.error(err)
            //context.client.setQueryData(['rooms'], onMutateResult?.prevRooms)
        },
        onSettled: (data, err, variables, onMutateResult, context) => {
            context.client.invalidateQueries({ queryKey: ['rooms'] })
        }
    })

    const deleteRoomQuery = useMutation({
        mutationKey: ['deleteRoom'],
        mutationFn: (room_id: string) => deleteRoom(room_id),
        onMutate: async (room_id, context) => {
            await context.client.cancelQueries({ queryKey: ['rooms', 'owned'] })

            const prevRooms = context.client.getQueryData(['rooms', 'owned'])

            context.client.setQueryData(['rooms', 'owned'], (old: Room[]) => old.filter((r: Room) => r.id !== room_id))

            return { prevRooms }
        },
        onError: (err, variables, onMutateResult, context) => {
            context.client.setQueryData(['rooms', 'owned'], onMutateResult?.prevRooms)
        },
        onSettled: (data, err, variables, onMutateResult, context) => {
            context.client.invalidateQueries({ queryKey: ['rooms', 'owned'] })
        }
    })

    const joinRoomQuery = useMutation({
        mutationKey: ['joinRoom'],
        mutationFn: (join_code: string) => joinRoom(join_code),
        //onMutate: ,
        onError: (err, newRoom, onMutateResult, context) => {
            console.error(err)
            //context.client.setQueryData([''], )
        },
        onSettled: (data, err, variables, onMutateResult, context) => {
            context.client.invalidateQueries({ queryKey: ['rooms', 'joined'] })
        }
    })

    const leaveRoomQuery = useMutation({
        mutationKey: ['leaveRoom'],
        mutationFn: (room_id: string) => leaveRoom(room_id),
        onMutate: async (room_id, context) => {
            await context.client.cancelQueries({ queryKey: ['rooms', 'owned'] })

            const prevRooms = context.client.getQueryData(['rooms', 'owned'])

            context.client.setQueryData(['rooms', 'joined'], (old: Room[]) => old.filter((r: Room) => r.id !== room_id))

            return { prevRooms }
        },
        onError: (err, variables, onMutateResult, context) => {
            context.client.setQueryData(['rooms', 'joined'], onMutateResult?.prevRooms)
        },
        onSettled: (data, err, variables, onMutateResult, context) => {
            context.client.invalidateQueries({ queryKey: ['rooms', 'joined'] })
        }
    })

    const handleDeleteRoom = (room_id: string) => {
        deleteRoomQuery.mutate(room_id)
    }

    const handleLeaveRoom = (room_id: string) => {
        leaveRoomQuery.mutate(room_id)
    }

    const handleCreateRoom = () => {
        setShowCreateRoom(true)
        setRoomName("")
        setRoomDescription("")
    }

    const handleCreateRoomSubmit = () => {
        setShowCreateRoom(false)
        addRoomQuery.mutate({ name: roomName, description: roomDescription })
    }

    const handleJoinRoom = () => {
        setShowJoinRoom(true)
    }

    const handleJoinRoomSubmit = () => {
        setShowJoinRoom(false)
        joinRoomQuery.mutate(joinCode)
    }

    const handleQRScan = () => {
        console.log("Open QR scanner")
    }

    if (ownedRoomsQuery.isPending || joinedRoomsQuery.isPending) {
        return (
            <div className="flex items-center justify-center min-h-[400px]">
                <div className="text-center">
                    <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
                    <p className="text-gray-600">Loading your rooms...</p>
                </div>
            </div>
        )
    }

    if (ownedRoomsQuery.isError || joinedRoomsQuery.isError) {
        return (
            <div className="flex items-center justify-center min-h-[400px]">
                <div className="text-center bg-red-50 border border-red-200 rounded-lg p-8 max-w-md">
                    <h2 className="text-xl font-semibold text-red-700 mb-2">Error Loading Rooms</h2>
                    <p className="text-red-600">
                        {(() => {
                            if (ownedRoomsQuery.isError && joinedRoomsQuery.isError) {
                                return ownedRoomsQuery.error.message + " and " + joinedRoomsQuery.error.message
                            } else if (ownedRoomsQuery.isError) {
                                return ownedRoomsQuery.error.message
                            } else {
                                return joinedRoomsQuery.error?.message
                            }
                        })()}
                    </p>
                </div>
            </div>
        )
    }

    const rooms = activeTab === "owned" ? ownedRoomsQuery.data : joinedRoomsQuery.data

    return (
        <div className="max-w-5xl mx-auto p-6">
            <div className="mb-8">
                <h1 className="text-3xl font-bold text-gray-900 mb-2">
                    Welcome back, {user != null ? user.username : "unknown"}!
                </h1>
                <p className="text-gray-600">Manage your learning rooms</p>
            </div>
            <div className="mb-6 flex justify-between">
                <Tabs value={activeTab} onValueChange={(value: string) => setActiveTab(value as "owned" | "joined")}>
                    <TabsList>
                        <TabsTrigger value="owned">My Rooms</TabsTrigger>
                        <TabsTrigger value="joined">Joined Rooms</TabsTrigger>
                    </TabsList>
                </Tabs>

                {activeTab === "owned" ? (
                    <Button onClick={handleCreateRoom} className="gap-2">
                        <Plus className="h-4 w-4" />
                        Create Room
                    </Button>
                    ) : (
                    <Button onClick={handleJoinRoom} className="gap-2">
                        <UserPlus className="h-4 w-4" />
                        Join Room
                    </Button>
                )}
            </div>

            <div className="space-y-3">
                {rooms && rooms.length > 0 ? (
                    rooms.map((room, i) => (
                        <RoomCard
                            key={i}
                            id={room.id}
                            roomName={room.name}
                            joinCode={room.join_code || "ABC123"}
                            isOwned={activeTab === "owned"}
                            onDelete={handleDeleteRoom}
                            onLeave={handleLeaveRoom}
                        />
                    ))
                ): (
                    <div className="text-center py-16 bg-gray-50 rounded-lg border-2 border-dashed border-gray-300">
                        <p className="text-gray-500 text-lg">
                            {activeTab === "owned" ? "You haven't created any rooms yet" : "You haven't joined any rooms yet"}
                        </p>
                        <p className="text-gray-400 text-sm mt-2">
                            {activeTab === "owned" ? "Create your first room to get started" : "Ask your instructor for a join code"}
                        </p>
                    </div>
                )}
            </div>

            <Dialog open={showCreateRoom} onOpenChange={setShowCreateRoom}>
                <DialogContent className="sm:max-w-md">
                <DialogHeader>
                    <DialogTitle>Create New Room</DialogTitle>
                    <DialogDescription>
                    Add a new learning room for your students to join
                    </DialogDescription>
                </DialogHeader>
                <div className="space-y-4 py-4">
                    <div className="space-y-2">
                    <Label htmlFor="room-name">Room Name *</Label>
                    <Input
                        id="room-name"
                        placeholder="e.g., Mathematics 101"
                        value={roomName}
                        onChange={(e) => setRoomName(e.target.value)}
                    />
                    </div>
                    <div className="space-y-2">
                    <Label htmlFor="room-description">Description (optional)</Label>
                    <Textarea
                        id="room-description"
                        placeholder="Add details about this room..."
                        value={roomDescription}
                        onChange={(e) => setRoomDescription(e.target.value)}
                        rows={3}
                    />
                    </div>
                </div>
                <DialogFooter>
                    <Button
                    variant="outline"
                    onClick={() => setShowCreateRoom(false)}
                    >
                    Cancel
                    </Button>
                    <Button
                    onClick={handleCreateRoomSubmit}
                    disabled={!roomName.trim()}
                    >
                    Create Room
                    </Button>
                </DialogFooter>
                </DialogContent>
            </Dialog>

            <Dialog open={showJoinRoom} onOpenChange={setShowJoinRoom}>
                <DialogContent className="sm:max-w-md">
                <DialogHeader>
                    <DialogTitle>Join a Room</DialogTitle>
                    <DialogDescription>
                    Enter a join code or scan a QR code to join a room
                    </DialogDescription>
                </DialogHeader>
                
                <Tabs value={joinMethod} onValueChange={(value) => setJoinMethod(value as "code" | "qr")}>
                    <TabsList className="grid w-full grid-cols-2">
                    <TabsTrigger value="code">Enter Code</TabsTrigger>
                    <TabsTrigger value="qr">Scan QR Code</TabsTrigger>
                    </TabsList>
                </Tabs>

                <div className="py-4">
                    {joinMethod === "code" ? (
                    <div className="space-y-2">
                        <Label htmlFor="join-code">Join Code</Label>
                        <Input
                        id="join-code"
                        placeholder="e.g., MATH101"
                        value={joinCode}
                        onChange={(e) => setJoinCode(e.target.value.toUpperCase())}
                        className="text-center text-lg font-mono tracking-wider"
                        />
                        <p className="text-sm text-gray-500">
                        Ask your instructor for the room join code
                        </p>
                    </div>
                    ) : (
                    <div className="space-y-4">
                        <div className="flex flex-col items-center justify-center p-8 border-2 border-dashed border-gray-300 rounded-lg bg-gray-50">
                        <Camera className="h-16 w-16 text-gray-400 mb-4" />
                        <p className="text-sm text-gray-600 text-center mb-4">
                            Position the QR code in front of your camera
                        </p>
                        <Button onClick={handleQRScan} variant="outline">
                            <Camera className="h-4 w-4 mr-2" />
                            Open Camera
                        </Button>
                        </div>
                        <p className="text-xs text-gray-500 text-center">
                        Your instructor can display a QR code for quick access
                        </p>
                    </div>
                    )}
                </div>

                <DialogFooter>
                    <Button
                    variant="outline"
                    onClick={() => {
                        setShowJoinRoom(false)
                        setJoinCode("")
                        setJoinMethod("code")
                    }}
                    >
                    Cancel
                    </Button>
                    {joinMethod === "code" && (
                    <Button
                        onClick={handleJoinRoomSubmit}
                        disabled={!joinCode.trim()}
                    >
                        Join Room
                    </Button>
                    )}
                </DialogFooter>
                </DialogContent>
            </Dialog>
        </div>
    )
}