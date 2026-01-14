import { RoomItem } from "@/components/room-item"
import { Button } from "@/components/ui/button"
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs"
import type { Room } from "@/lib/api/models/rooms"
import { useAuth } from "@/lib/context/AuthContext"
import useRoomsApi from "@/lib/hooks/useRoomsApi"
import { useMutation, useQuery } from "@tanstack/react-query"
import { useState } from "react"
import { Copy, QrCode, Trash2, LogOut, Plus, UserPlus } from "lucide-react"
import { data } from "react-router-dom"

export default function HomePage() {
    const { user } = useAuth()
    const { getJoinedRooms, getOwnedRooms, createRoom, deleteRoom } = useRoomsApi()
    const [activeTab, setActiveTab] = useState<"owned" | "joined">("owned")
    const [showCreateRoom, setShowCreateRoom] = useState(false)
    const [showJoinRoom, setShowJoinRoom] = useState(false)
    const [joinMethod, setJoinMethod] = useState<"code" | "qr">("code")
    const [joinCode, setJoinCode] = useState("")
    const [roomName, setRoomName] = useState("")
    const [roomDescription, setRoomDescription] = useState("")


    const ownedRoomsQuery = useQuery({
        queryKey: ['rooms'],
        queryFn: getOwnedRooms
    })

    const joinedRoomsQuery = useQuery({
        queryKey: ['rooms'],
        queryFn: getJoinedRooms
    })

    const handleDeleteRoom = (room_id: string) => {
        useMutation({
            mutationKey: ['deleteRoom'],
            mutationFn: deleteRoom,
            onMutate: async (variables, context) => {
                await context.client.cancelQueries({ queryKey: ['rooms'] })

                const prevRooms = context.client.getQueryData(['rooms'])

                context.client.setQueryData(['rooms'], (old: Room[]) => old.filter((r: Room) => r.id !== room_id))

                return { prevRooms }
            },
            onError: (err, variables, onMutateResult, context) => {
                context.client.setQueryData(['rooms'], onMutateResult?.prevRooms)
            },
            onSettled: (data, err, variables, onMutateResult, context) => {
                context.client.invalidateQueries({ queryKey: ['rooms'] })
            }
        })
    }

    const handleCreateRoom = () => {
        setShowCreateRoom(true)
        setShowCreateRoom(false)
        setRoomName("")
        setRoomDescription("")
    }

    const handleCreateRoomSubmit = (room_id: string) => {
        useMutation({
            mutationKey: ['addRoom'],
            mutationFn: createRoom,
            onMutate: async (newRoom, context) => {
                await context.client.cancelQueries({ queryKey: ['rooms'] })
                
                const prevRooms = context.client.getQueryData(['rooms'])

                context.client.setQueryData(['rooms'], (old: Room[]) => [...old, newRoom])

                return { prevRooms }
            },
            onError: (err, newRoom, onMutateResult, context) => {
                context.client.setQueryData(['rooms'], onMutateResult?.prevRooms)
            },
            onSettled: (data, err, variables, onMutateResult, context) => {
                context.client.invalidateQueries({ queryKey: ['rooms'] })
            }
        })
    }

    const handleJoinRoom = () => {
        setShowJoinRoom(true)
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

    return (
        <div className="max-w-5xl mx-auto p-6">
            <div className="mb-8">
                <h1 className="text-3xl font-bold text-gray-900 mb-2">
                    Welcome back, {user != null ? user.username : "unknown"}!
                </h1>
                <p className="text-gray-600">Manage your learning rooms</p>
            </div>
            <div className="mb-6">
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
                {activeTab === "owned" && ownedRoomsQuery.data && ownedRoomsQuery.data.length > 0 ? (
                    ownedRoomsQuery.data.map((room, i) => (
                        <RoomItem
                            key={i}
                            id={room.id}
                            roomName={room.name}
                            joinCode={room.join_code || "ABC123"}
                            onDelete={handleDeleteRoom}
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
        </div>
    )
}