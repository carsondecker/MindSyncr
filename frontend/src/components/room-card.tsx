import { Button } from "@/components/ui/button"
import { Link } from "react-router-dom"
import { Copy, LogOut, QrCode, Trash2 } from "lucide-react"
import { Card, CardDescription, CardFooter, CardHeader, CardTitle } from "./ui/card"
import { useState } from "react"
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from "./ui/dialog"
import { QRCodeCanvas } from 'qrcode.react'

type RoomCardProps = {
    id: string
    roomName: string
    joinCode: string
    isOwned: boolean
    onDelete?: (room_id: string) => void
    onLeave?: (room_id: string) => void
}

export function RoomCard({ id, roomName, joinCode, isOwned, onDelete, onLeave }: RoomCardProps) {
  const [showJoinCode, setShowJoinCode] = useState(false)
  const [copied, setCopied] = useState(false)

  const getLink = () => `${window.location.origin}/join/${joinCode}`

  const handleCopyLink = () => {
    const shareableLink = getLink()
    navigator.clipboard.writeText(shareableLink)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  const handleShowJoinCode = () => {
    setShowJoinCode(true)
  }

  return (
    <>
      <Card className="hover:shadow-lg transition-shadow border-l-4 border-l-blue-500">
        <CardHeader>
          <CardTitle>
            <Link to={`/rooms/${id}`} className="hover:text-blue-600 transition-colors">
              {roomName}
            </Link>
          </CardTitle>
          <CardDescription>Room ID: {id}</CardDescription>
        </CardHeader>
        <CardFooter className="flex gap-2 justify-end">
          {isOwned ? (
            <>
              <Button
                variant="outline"
                size="sm"
                onClick={handleCopyLink}
                className="gap-2"
              >
                <Copy className="h-4 w-4" />
                {copied ? "Copied!" : "Copy Link"}
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={handleShowJoinCode}
                className="gap-2"
              >
                <QrCode className="h-4 w-4" />
                Join Code
              </Button>
              <Button
                variant="destructive"
                size="sm"
                onClick={() => onDelete?.(id)}
                className="gap-2"
              >
                <Trash2 className="h-4 w-4" />
                Delete
              </Button>
            </>
          ) : (
            <Button
              variant="outline"
              size="sm"
              onClick={() => onLeave?.(id)}
              className="gap-2 hover:bg-red-50 hover:text-red-600 hover:border-red-300"
            >
              <LogOut className="h-4 w-4" />
              Leave
            </Button>
          )}
        </CardFooter>
      </Card>

      <Dialog open={showJoinCode} onOpenChange={setShowJoinCode}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Join Code for {roomName}</DialogTitle>
            <DialogDescription>
              Share this code with others to let them join your room
            </DialogDescription>
          </DialogHeader>
          <div className="flex items-center justify-center p-4">
            <div className="text-center">
              <p className="text-sm text-gray-500 mb-2">
                Students can use this code to join
              </p>
              <div className="inline-block bg-gradient-to-br from-blue-50 to-indigo-50 border-2 border-blue-300 rounded-lg px-8 py-6">
                <p className="text-4xl font-bold tracking-wider text-blue-700 font-mono">
                  {joinCode}
                </p>
              </div>
              <div className="relative flex py-5 items-center">
                <div className="flex-grow border-t border-gray-400"></div>
                <span className="flex-shrink mx-4 text-sm text-gray-400">OR</span>
                <div className="flex-grow border-t border-gray-400"></div>
              </div>
              <p className="text-sm text-gray-500 mb-2">
                Students can scan this code to join
              </p>
              <div className="flex items-center justify-center">
                <QRCodeCanvas
                  value={getLink()}
                  size={256}
                  style={{ width: "75%", height: "75%" }}
                />
              </div>
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </>
  )
}