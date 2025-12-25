import { useAuth } from "@/lib/context/AuthContext"

export default function HomePage() {
    const { user } = useAuth()

    return (
        <h1>
            Welcome back, {user != null ? user.username : "unknown"}
        </h1>
    )
}