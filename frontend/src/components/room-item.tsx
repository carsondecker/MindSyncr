import { Button } from "@/components/ui/button"
import { Link } from "react-router-dom"
import { Item, ItemActions, ItemContent, ItemTitle } from "./ui/item"

type RoomItemProps = {
    id: string
    roomName: string
    joinCode: string
    onDelete: (room_id: string) => void
}

export function RoomItem({id, roomName, joinCode}: RoomItemProps) {
  return (
    <Item variant="outline">
      <ItemContent>
        <ItemTitle>
          <Link to={`/rooms/${id}`}>{roomName}</Link>
        </ItemTitle>
      </ItemContent>
      <ItemActions>
        <Button variant="destructive" size="sm">
          Delete
        </Button>
      </ItemActions>
    </Item>
  )
}
