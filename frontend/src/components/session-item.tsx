import { Button } from "@/components/ui/button"
import { Link } from "react-router-dom"
import { Item, ItemActions, ItemContent, ItemTitle } from "./ui/item"
import { useApi } from "@/lib/hooks/useApi"
import { useEffect } from "react"
import { deleteSession, endSession } from "@/lib/api/sessions"

type SessionItemProps = {
    id: string
    room_id: string
    sessionName: string
    is_active: boolean
    owner_id: string
    ended_at: Date | null
    is_owner: boolean
    is_member: boolean
    removeItem: (session_id: string) => void
    endItem: (session_id: string) => void
}

export function SessionItem({ id, room_id, sessionName, is_active, owner_id, is_owner, is_member, ended_at, removeItem, endItem }: SessionItemProps) {
  const { run, loading, error } = useApi()
  
  const deleteSelf = async () => {
    try {
      await run(() => deleteSession(id))
      removeItem(id)
    } catch (err) {

    }
  }

  const endSelf = async () => {
    try {
      await run(() => endSession(id))
      endItem(id)
    } catch (err) {

    }
  }

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
                <Link to={`sessions/${id}`}>Open</Link>
              </Button>
              <Button variant="destructive" size="sm" onClick={deleteSelf}>
                  Delete
              </Button>
              {is_active && !ended_at && (
                  <Button size="sm" onClick={endSelf}>
                      End Session
                  </Button>
              )}
          </>
        )}
        {!is_owner && (
          <>
            {!is_member && (
              <Button variant="outline" size="sm">
                <Link to={`sessions/${id}`}>Rejoin Session</Link>
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
