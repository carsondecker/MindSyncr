import { Button } from "@/components/ui/button"
import { Link } from "react-router-dom"
import { Item, ItemActions, ItemContent, ItemTitle } from "./ui/item"
import useSessionsApi from "@/lib/hooks/useSessionsApi"

type SessionItemProps = {
    session_id: string
    room_id: string
    sessionName: string
    is_active: boolean
    owner_id: string
    ended_at: Date | null
    is_owner: boolean
    is_member: boolean
    deleteItem: (session_id: string) => void
    endItem: (session_id: string) => void
}

export function SessionItem({ session_id, room_id, sessionName, is_active, owner_id, is_owner, is_member, ended_at, deleteItem, endItem }: SessionItemProps) {
  return (
    <Item variant="outline">
      <ItemContent>
        <ItemTitle>
          {sessionName}
        </ItemTitle>
      </ItemContent>
      <ItemActions>
        {is_owner && (
          <>
              <Button size="sm">
                <Link to={`/sessions/${session_id}`}>Open</Link>
              </Button>
              <Button variant="destructive" size="sm" onClick={() => deleteItem}>
                  Delete
              </Button>
              {is_active && !ended_at && (
                  <Button size="sm" onClick={() => endItem}>
                      End Session
                  </Button>
              )}
          </>
        )}
        {!is_owner && (
          <>
            {!is_member && (
              <Button variant="outline" size="sm">
                <Link to={`/sessions/${session_id}`}>Rejoin Session</Link>
              </Button>
            )}
            {is_member && (
              <Button variant="outline" size="sm">
                  Join Session
              </Button>
            )}
          </>
        )}
      </ItemActions>
    </Item>
  )
}
