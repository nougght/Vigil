// import { Link, useNavigate } from "react-router-dom"
import {useState} from "react";
// import { GroupCard } from "../components/groupCard";
import styles from "./groupsPage.module.css"
import TabBar from "../../../shared/ui/src/components/TabBar.tsx";
import {Checkbox, Dialog} from "radix-ui";
import {CheckIcon} from "@radix-ui/react-icons";
import {useGroups} from "../hooks/useGetGroups.ts";
import {GroupCard} from "../components/groupCard.tsx";
import {useCreateGroup} from "../hooks/useCreateGroup.ts";
import {useAgents} from "../hooks/useGetAgents.ts";
import type {Agent} from "../domain/agent";
import {TextInput} from "@monitoring-system/shared/src/components/TextInput"
import {useNavigate} from "react-router-dom";

interface GroupInfo {
  Name: string
  Description?: string
  Agents?: string[]
}

const AgentsPicker = ({agents, checked, onChange}: {
  agents: Agent[],
  checked: Set<string>,
  onChange: (checked: Set<string>) => void
}) => {
  const toggleChecked = (id: string) => {
    const newSet = new Set(checked)
    if (newSet.has(id)) {
      newSet.delete(id)
    } else {
      newSet.add(id)
    }
    onChange(newSet)
  }

  return (
    <div>
      {
        agents && agents.map(
          a => (<div className={styles.checkRow}>
            <Checkbox.Root className={styles.checkboxRoot} checked={checked.has(a.id)}
                           onCheckedChange={() => toggleChecked(a.id)}
                           id={a.id}>
              <Checkbox.Indicator className={styles.checkboxInd}>
                <CheckIcon/>
              </Checkbox.Indicator>
            </Checkbox.Root>
            <label className="Label" htmlFor="c1">
              {a.name}
            </label>
          </div>)
        )
      }
    </div>
  )
}


interface tab {
  title: string
  content: React.ReactNode
  nextTitle: string
}

export const GroupsPage = () => {
  const [activeTab, setActiveTab] = useState(0)
  const [groupInfo, setGroupInfo] = useState<GroupInfo>({
    Name: ""
  })
  const [agentsChecked, setAgentsChecked] = useState<Set<string>>(new Set())
  const [nameWarning, setNameWarning] = useState<string>()

  const navigate = useNavigate()

  const {
    data: groups,
    isPending: isGroupsPending,
    isError: _isGroupsError,
    error: _groupsError,
    isFetching: _isGroupsFetching,
  } = useGroups();

  const {
    mutate: createGroup,
    isPending: isCreatePending,
    isError: _isCreateError,
    isSuccess: _isCreateSuccess,
    error: _createError,
  } = useCreateGroup()

  const {
    data: agents,
    isPending: _isAgentsPending,
    isError: _isAgentsError,
    error: _agentsError,
    isFetching: _isAgentsFetching,

  } = useAgents();


  const handleNext = () => {
    if (activeTab == 0) {
      if (!groupInfo.Name) {
        setNameWarning("Название группы не может быть пустым")
        return
      }
      setActiveTab(1);
    }
    createGroup(
      {
        name: groupInfo.Name,
        description: groupInfo.Description,
        agentIDs: Array.from(agentsChecked)
      },
    )
  }

  const tabs: tab[] = [
    {
      title: "Основное",
      nextTitle: "Далее",
      content:
        <div>
          <label className={styles.formLabel} htmlFor="name">Название</label>
          <TextInput value={groupInfo?.Name}
                     placeholder="Введите название группы"
                     warning={nameWarning}
                     onChange={value => {
                       if (value.length > 0)
                         setNameWarning(undefined);
                       setGroupInfo({
                         ...groupInfo,
                         Name: value
                       })
                     }
                     }/>
          <br/>
          <label className={styles.formLabel} htmlFor="description">Описание</label>
          <TextInput
            value={groupInfo?.Description}
            placeholder="Введите описание группы"
            multiline={true}
            onChange={v =>
              setGroupInfo({
                ...groupInfo,
                Description: v
              })}/>

          {/*{nameWarning && <p className={styles.warning}>{nameWarning}</p>}*/}


        </div>
    },
    {
      title: "Добавить агентов",
      nextTitle: "Создать",
      content:
        <div>
          <AgentsPicker agents={agents ?? []}
                        checked={agentsChecked}
                        onChange={setAgentsChecked}/>

          {_isAgentsPending && <a>Загрузка...</a>}
        </div>
    }
  ]

  const sidePanel = (
    <Dialog.Root>
      <Dialog.Trigger className={styles.createButton}>Новая группа</Dialog.Trigger>
      <Dialog.Portal>
        <Dialog.Overlay className={styles.overlay}/>
        <Dialog.Content className={styles.content}>
          <div className={styles.sidePanelContainer}>
            <h2>Создание группы</h2>
            <TabBar
              className={styles.tabBar}
              tabs={tabs.map((tab) => tab.title)}
              onSwitch={setActiveTab}
              activeTab={activeTab}/>

            <div className={styles.sidePanelContent}>
              {tabs[activeTab].content}
            </div>

            <div className={styles.bottomArea}>
              {isCreatePending && <a>Создание...</a>}
              <button onClick={handleNext}>{tabs[activeTab].nextTitle}</button>
            </div>
          </div>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  )


  return (
    <div className={styles.groupsPage}>
      <h1>Группы</h1>
      <main>

        {isGroupsPending && <a>Загрузка</a>}
        <div className={styles.topNav}>
          <div></div>
          {sidePanel}
        </div>
        <div className={styles.groupCardsContainer}>{
          groups && groups.map((g) =>
            <GroupCard key={g.id} group={g}
                       onClick={() => {
                         navigate(g.id)
                       }}/>
          )
        }

        </div>


      </main>
      {
        _isCreateError &&
          <div>
              <p>{_createError.message}</p>
          </div>
      }
    </div>
  )
}

