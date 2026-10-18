import styles from "./GroupPage.module.css"
import {useAgents} from "../hooks/useGetAgents.ts";
import {useGroup} from "../hooks/useGetGroups.ts";
import {type ReactNode, useState} from "react";
import {useNavigate, useParams} from "react-router-dom";
import TabBar from "../../../shared/ui/src/components/TabBar.tsx";
import {Table} from "@monitoring-system/shared/src/components/Table.tsx";


interface AgentRow {
  id: string;
  name: string
  lastSeenAt?: Date
}

interface ColumnInfo<T> {
  key: keyof T;
  header: string;
  sortable: boolean;
  fieldRenderFunc?: (row: T) => ReactNode
}

const Columns: Record<string, ColumnInfo<AgentRow>> = {
  "name": {
    key: "name",
    header: "Агент",
    sortable: true,
  },
  "lastSeenAt": {
    key: "lastSeenAt",
    header: "Последнее подключение",
    sortable: true,
    fieldRenderFunc: row => {
      return <span>{row["lastSeenAt"]?.toDateString() ?? ""}</span>
    }
  }
}

interface tab {
  title: string
  content: ReactNode
}

export const GroupPage = () => {
  const [selectedIDs, setSelected] = useState<Set<string>>(new Set<string>());
  const [sortBy, setSortBy] = useState("name");
  const [sortOrder, setSortOrder] = useState<"asc" | "desc">("asc");
  const {id} = useParams()

  const [activeTab, setActiveTab] = useState(0)

  const navigate = useNavigate();

  const {
    data: agents,
    isPending: isAgentsPending,
    isError: _isAgentsError,
    error: _agentsError,
    isFetching: _isAgentsFetching,
  } = useAgents()

  const {
    data: group,
    isPending: _isGroupPending,
    isError: _isGroupError,
    error: _groupError,
    isFetching: _isGroupFetching,
  } = useGroup(id ?? "")

  const tabs: tab[] = [
    {
      title: "Обзор",
      content: <></>
    },
    {
      title: "Агенты",
      content:
        <div className={styles.agentsContainer}>
          {isAgentsPending && "Загрузка..."}
          <Table<AgentRow>
            rows={agents?.map(a => {
              return {
                id: a.id,
                name: a.name,
                lastSeenAt: a.lastSeenAt,
              }
            }) ?? []}
            onRowClick={(id) => {
              navigate(`../agents/${id}`)
            }}
            columns={(Object.keys(Columns) as Array<keyof AgentRow>).map(key => {
              return Columns[key]
            })}
            selected={selectedIDs}
            onSelect={(ids, all) => {
              if (all === true) {
                if (selectedIDs.size != agents?.length) {
                  setSelected(agents ? new Set(agents.map(a => a.id) ?? []) : selectedIDs)
                  return
                }
                setSelected(new Set())
                return
              }

              const newSelected = new Set(selectedIDs)
              console.log(ids)
              ids.forEach(id => {
                console.log(id)
                if (newSelected.has(id)) {
                  newSelected.delete(id)
                  return
                }
                newSelected.add(id)
              })
              console.log(newSelected)
              setSelected(newSelected)
            }}
            sortColumn={sortBy}
            sortOrder={sortOrder}
            onSort={(column) => {
              if (column == sortBy) {
                setSortOrder(
                  sortOrder == "desc" ? "asc" : "desc"
                )
                return
              }
              setSortBy(column)
              setSortOrder("asc")
            }

            }
          />
        </div>
    }
  ]

  return (
    <div className={styles.GroupPage}>
      <div>
        <h2>{`Группа ${group?.name}`}</h2>
        <TabBar
          className={styles.tabBar}
          tabs={tabs.map((tab) => tab.title)}
          onSwitch={setActiveTab}
          activeTab={activeTab}/>
      </div>
      <div>
        {tabs[activeTab].content}
      </div>
    </div>)
}