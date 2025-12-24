import { useEffect } from "react"
import { handleLogin } from "./lib/api/auth"
import LoginPage from "./pages/LoginPage"

function App() {
  useEffect(() => {
    async function loginTest() {
      await handleLogin({ email: "test@gmail.com", password: "test" })
    }
    loginTest()
  }, [])
  

  return (
    <>
      <LoginPage />
    </> 
  )
}

export default App