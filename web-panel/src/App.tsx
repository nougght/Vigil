import './App.css'
import { BrowserRouter, Navigate, Route, Routes, Outlet, useLocation } from 'react-router-dom'
import { AgentsPage } from './pages/AgentsPage'
import { AgentPage } from './pages/AgentPage'
import { NewAgentPage } from './pages/NewAgentPage'
import { SideBar, type SideBarData } from './components/sideBar'
import { useEffect, useState } from 'react'
import dashIcon from "./assets/dashboard.svg"
import agentsIcon from "./assets/cpu.svg"
import { ClientMessageTypeAgentDetailed, FillMetricsFromSeries, MessageTypeSeries, type AgentDetailedMessage, type ClientMessage, type Message, type Metrics } from './domain/metrics'
// import type { SeriesDTO } from './domain/metrics'
const sideBarData: SideBarData = {
    iconSrc: "",
    title: "Vigil",
    items: [
        { id: "dashboard", title: "Обзор", iconSrc: dashIcon, path: "/dashboard", countLabel: 0 },
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
    const [metrics, setMetrics] = useState<Metrics | undefined>()

    useEffect(() => {
        const socket = new WebSocket("ws://monitoring.nought.ru/api/v1/ws");
        socket.addEventListener("open", () => {
            console.log("start")
        });

        socket.onclose = () => {
            console.log("connection closed")
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
        if (wsSocket == null || wsSocket.readyState == wsSocket.CONNECTING) {
            console.log("ws socket is null")
            return
        }
        wsSocket.send(JSON.stringify(sendMessage))
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
                        <Route path="/dashboard" element={<p>not implemented</p>} />
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

