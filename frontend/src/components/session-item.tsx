import { Button } from "@/components/ui/button"
import { Link } from "react-router-dom"
import { Item, ItemActions, ItemContent, ItemTitle } from "./ui/item"
import { useAuth } from "@/lib/context/AuthContext"

type SessionItemProps = {
    room_id: string
    sessionName: string
    is_active: boolean
    owner_id: string
    ended_at: Date | null
}

export function SessionItem({ room_id, sessionName, is_active, owner_id, ended_at }: SessionItemProps) {
    const { user, loading  } = useAuth()

    return (
    <Item variant="outline">
      <ItemContent>
        <ItemTitle>
          <Link to={`/rooms/${room_id}`}>{sessionName}</Link>
        </ItemTitle>
      </ItemContent>
      <ItemActions>
        {user?.id == owner_id && (
            <>
                <Button variant="destructive" size="sm">
                    Delete
                </Button>
                {is_active && !ended_at && (
                    <Button>
                        End Session
                    </Button>
                )}
            </>
        )}
         {user?.id != owner_id && (
            <Button variant="outline" size="sm">
                Join Session
            </Button>
        )}
      </ItemActions>
    </Item>
  )
}
