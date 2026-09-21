import { Link, useNavigate } from "react-router-dom"
import { useAgents } from "../hooks/useGetAgents";
import { useEffect, useState } from "react";
import { AgentCard } from "../components/agentCard";
import styles from "./agentsPage.module.css"



export const AgentsPage = () => {
    const [warning, setWarning] = useState<string | null>()
    const {
        data: agents,
        isPending: isAgentsPending,
        isError: _isAgentsError,
        error: _agentsError,
        isFetching: _isAgentsFetching,
    } = useAgents();
    let navigate = useNavigate()


    useEffect(() => {
        if (agents?.error != null) {
            setWarning(`ошибка:${agents?.error.status} ${agents?.error.message}`)
        }
    }, [agents]);

    if (isAgentsPending) {
        return <div>Загрузка...</div>;
    }
    return (
        <div className={styles.agentsPage}>
            <h1>Агенты</h1>
            <main>
                <div className={styles.agentCardsContainer}>
                    {agents?.agents != null && agents?.agents.length > 0 &&
                        agents?.agents?.filter((a) => a.status != null)
                            .sort((a, b) => {
                                if (a.isOnline == b.isOnline)
                                    return 0
                                if (a.isOnline) {
                                    return -1
                                }
                                return 1
                            })
                            .map((agent) =>
                                // {agent.status != null &&
                                <div key={agent.id}>
                                    <AgentCard agent={agent} onClick={(id) => { navigate(id) }} />
                                </div>
                                // }
                            )
                    }
                </div>
                <div className={styles.bottomArea}>
                    <button className={styles.addAgentButton}>
                        <Link to="./new">Добавить</Link>
                    </button>
                </div>
            </main>
            {
                warning != null &&
                <div>
                    <p>{warning}</p>
                </div>
            }
        </div>
    )
}