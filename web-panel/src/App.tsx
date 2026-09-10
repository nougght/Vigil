import './App.css'
import { BrowserRouter, Navigate, Route, Routes, Outlet, useLocation } from 'react-router-dom'
import { AgentsPage } from './pages/AgentsPage'
import { AgentPage } from './pages/AgentPage'
import { NewAgentPage } from './pages/NewAgentPage'
import { SideBar, type SideBarData } from './components/sideBar'
import { useEffect, useRef, useState } from 'react'
import dashIcon from "./assets/dashboard.svg"
import agentsIcon from "./assets/cpu.svg"
import { ClientMessageTypeAgentDetailed, FillMetricsFromSeries, MessageTypeSeries, type AgentDetailedMessage, type ClientMessage, type Message, type Metrics } from './domain/metrics'
import { OverviewPage } from './pages/OverviewPage'
// import type { SeriesDTO } from './domain/metrics'
const sideBarData: SideBarData = {
    iconSrc: "",
    title: "Vigil",
    items: [
        { id: "overview", title: "Обзор", iconSrc: dashIcon, path: "/overview", countLabel: 0 },
        { id: "agents", title: "Агенты", iconSrc: agentsIcon, path: "/agents", countLabel: 0 },
    ]

}
const AppLayout = () => {
    let location = useLocation()

    useEffect(
        () => {

        },
        [location]
    )
    return (
        <div style={{ display: `flex`, flexDirection: `row`, height: `100%`, alignItems: `stretch` }}>
            <SideBar
                data={sideBarData}
            />
            <div style={{ flexGrow: 1, overflow: `auto` }}><Outlet /></div>

        </div>
    )
}


function App() {
    const [sendMessage, setSendMessage] = useState<ClientMessage | null>(null)
    const [wsSocket, setWSSocket] = useState<WebSocket | null>(null)
    const [socketConnected, setConnected] = useState<boolean>(false)
    const [metrics, setMetrics] = useState<Metrics | undefined>()

    const sendMessageRef = useRef(sendMessage);

    const sendMsg = () => {
         const currentMsg = sendMessageRef.current; 
        if (currentMsg != null) {
            console.log("send message ", currentMsg)
            wsSocket?.send(JSON.stringify(currentMsg))
        }
    }
    useEffect(() => {
        const socket = new WebSocket("ws://monitoring.nought.ru/api/v1/ws");
        socket.addEventListener("open", () => {
            console.log("start")
            setConnected(true)
        });

        socket.onclose = () => {
            console.log("connection closed")
            setConnected(false)
        };

        socket.addEventListener("message", (event) => {
            // console.log("Message from server ", event.data);
            const msg = JSON.parse(event.data) as Message;
            if (msg.type == MessageTypeSeries) {
                console.log("series message received", msg)
                const series = msg.payload
                if (series != undefined) {
                    setMetrics((prev) => FillMetricsFromSeries(series, prev ?? {}))
                }
            }
        });

        setWSSocket(socket)
    }, []);

    useEffect(() => {
        sendMsg()
    }, [socketConnected])

    useEffect(() => {
        if (sendMessage == null)
            return
        sendMessageRef.current = sendMessage
        if (wsSocket == null || wsSocket.readyState != wsSocket.OPEN) {
            console.log("ws socket not ready: ", wsSocket?.readyState)
            return
        }
        console.log("ws message sent ", sendMessage)
        sendMsg()
    }, [sendMessage]);

    const handleLocationChange = (path: string) => {
        const splitted = path.split("/")
        if (splitted[1] == "agents" && splitted.length == 3) {
            console.log("agent id = ", splitted[2])
            const p: AgentDetailedMessage = {
                agents: [splitted[2]]
            }
            const msg: ClientMessage = {
                type: ClientMessageTypeAgentDetailed,
                payload: p
            }
            setSendMessage(msg)
        }
    }

    return (
        <div className='app'>
            <BrowserRouter>
                <RouteChangeTracker
                    handler={handleLocationChange}
                />
                <Routes>
                    <Route element={<AppLayout />}>
                        <Route path="/" element={<Navigate to="/agents" replace />} />
                        <Route path="/overview" element={<OverviewPage />} />
                        <Route path="/agents" element={<AgentsPage />} />
                        <Route path="/agents/:id" element={<AgentPage
                            metrics={metrics} />} />
                        <Route path="/agents/new" element={<NewAgentPage />} />
                    </Route>
                    {/* <Route path="*" element={<NotFoundPage />} /> */}
                </Routes>
            </BrowserRouter>
        </div>
    )
}


function RouteChangeTracker({ handler }: { handler: (path: string) => void }) {
    const location = useLocation();

    useEffect(() => {
        console.log('Route changed to:', location.pathname);
        handler(location.pathname)

    }, [location]);

    return null; // This component doesn't need to render anything visual
}

export default App

