import {Link, useNavigate} from "react-router-dom"
import {useAgents} from "../hooks/useGetAgents";
import {AgentCard} from "../components/agentCard";
import styles from "./agentsPage.module.css"
import {useState} from "react";
import {TextInput} from "../../../shared/ui/src/components/TextInput.tsx";


export const AgentsPage = () => {
  const [search, setSearch] = useState<string>()

  const {
    data: agents,
    isPending: isAgentsPending,
    isError: _isAgentsError,
    error: _agentsError,
    refetch: refetchAgents,
    isFetching: _isAgentsFetching,
  } = useAgents({
    search: search
  });

  const navigate = useNavigate()

  return (
    <div className={styles.agentsPage}>
      <h1>Агенты</h1>
      <main>
        <TextInput
          className={styles.searchField}
          onChange={(value) => {
            setSearch(value);
            if (!isAgentsPending) {
              refetchAgents();
            }
          }}
          value={search}
          placeholder={"Поиск"}
        />
        {isAgentsPending ? <div>Загрузка...</div> :
          <div className={styles.agentCardsContainer}>
            {agents && agents?.length > 0 &&
              agents.filter((a) => a.status != null)
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
                      <AgentCard agent={agent} onClick={(id) => {
                        navigate(id)
                      }}/>
                    </div>
                  // }
                )
            }
          </div>
        }
        <div className={styles.bottomArea}>
          <button className={styles.addAgentButton}>
            <Link to="./new">Добавить</Link>
          </button>
        </div>
      </main>
      {
        _isAgentsError &&
          <div>
              <p>ошибка: {_agentsError.message}</p>
          </div>
      }
    </div>
  )
}