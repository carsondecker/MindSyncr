import { Button } from "@/components/ui/button"
import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"

import { Link } from "react-router-dom"

type RoomCardProps = {
    id: string
    roomName: string
    roomDescription: string
    joinCode: string

}

export function RoomCard({id, roomName, roomDescription, joinCode}: RoomCardProps) {
  return (
    <Card className="w-full max-w-sm">
      <CardHeader>
        <CardTitle>
            <Link to="/" className="text-primary underline-offset-4 hover:underline">
                {roomName}
            </Link>
        </CardTitle>
        <CardDescription>
          {roomDescription}
        </CardDescription>
        <CardAction>
          <Button variant="destructive">Delete</Button>
        </CardAction>
      </CardHeader>
    </Card>
  )
}
