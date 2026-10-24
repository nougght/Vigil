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
import {type Agent, IsDefaultGroup} from "../domain/agent";
import {TextInput} from "@monitoring-system/shared/src/components/TextInput"
import {useNavigate} from "react-router-dom";
import {DropDownPicker} from "@monitoring-system/shared/src/components/DropDownPicker.tsx";
import type {Group} from "../domain/group.ts";
import {ToggleChips} from "@monitoring-system/shared/src/components/ToggleChips.tsx";
import {DefaultId} from "../domain/common.ts";

interface GroupInfo {
  Name: string
  Description?: string
  Agents?: string[]
}


const AgentsPicker = ({agents, groupNames, checked, onChange}: {
  agents: Agent[],
  groupNames: { [k: string]: string },
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
    <div className={styles.agentsPicker}>
      {
        agents && agents.map(
          a => (
            <div className={styles.checkRow}>
              <Checkbox.Root className={styles.checkboxRoot} checked={checked.has(a.id)}
                             onCheckedChange={() => toggleChecked(a.id)}
                             id={a.id}>
                <Checkbox.Indicator className={styles.checkboxInd}>
                  <CheckIcon/>
                </Checkbox.Indicator>
              </Checkbox.Root>
              <div className={styles.checkRowLabelCont}>
                <span>
                  {a.name}
                </span>
                <span className={styles.groupLabel}>
                  {IsDefaultGroup(a) ? "Без группы" : groupNames[a.groupID]}
                </span>
              </div>
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
  const [agentSearch, setAgentSearch] = useState<string>()

  const [agentsChecked, setAgentsChecked] = useState<Set<string>>(new Set())
  const [nameWarning, setNameWarning] = useState<string>()

  const [pickerGroupsChecked, setPickerGroupsChecked] = useState<string[]>([DefaultId])
  const [pickerModeToggle, setPickerModeToggle] = useState("without")

  const navigate = useNavigate()

  const {
    data: groups,
    isPending: isGroupsPending,
    isError: _isGroupsError,
    error: _groupsError,
    isFetching: _isGroupsFetching,
  } = useGroups();

  const pickerGroupItems = groups?.map((group: Group) => {
    return {
      value: group.id,
      label: group.name,
    }
  }) ?? []

  const {
    mutate: createGroup,
    isPending: isCreatePending,
    isError: _isCreateError,
    isSuccess: _isCreateSuccess,
    error: _createError,
  } = useCreateGroup()

  const groupNames = Object.fromEntries(groups?.map(g => [g.id, g.name]) ?? [])

  const {
    data: agents,
    isPending: _isAgentsPending,
    isError: _isAgentsError,
    error: _agentsError,
    refetch: refetchAgents,
    isFetching: _isAgentsFetching,
  } = useAgents({
    search: agentSearch,
    groupIds: pickerGroupsChecked,
  });


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

  const toggleGroupItems = [
    {
      value: "without",
      label: "Без группы"
    },
    {
      value: "select",
      label: "Выбрать группу"
    },
    {
      value: "all",
      label: "Все группы"
    },
  ]

  const tabs: tab[] = [
    {
      title: "Основное",
      nextTitle: "Далее",
      content:
        <div className={styles.sidePanelTabContent + " " + styles.sidePanelMainTabContent}>
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
        <div className={styles.sidePanelTabContent}>
          <div className={styles.sidePanelAgentsTop}>
            <ToggleChips
              className={styles.pickerGroupModeToggle}
              items={toggleGroupItems}
              multiple={false}
              defaultValue={pickerModeToggle}
              onChange={(value) => {
                setPickerModeToggle(value)
                if (value === "all") {
                  setPickerGroupsChecked([])
                }
                if (value === "select") {
                  setPickerGroupsChecked([])
                }
                if (value === "without") {
                  setPickerGroupsChecked([DefaultId])
                }
              }
              }
            />
            {
              pickerModeToggle === "select" &&
                <DropDownPicker items={pickerGroupItems}
                                multiple={true}
                                triggerPlaceholder={"Выберите группу"}
                                checked={pickerGroupsChecked}
                                onChange={setPickerGroupsChecked}
                                needSearch={true}/>
            }
            <TextInput
              className={styles.searchField}
              onChange={(value) => {
                setAgentSearch(value);
                if (!_isAgentsPending) {
                  refetchAgents();
                }
              }}
              value={agentSearch}
              placeholder={"Поиск"}
            />
          </div>
          <AgentsPicker agents={agents ?? []}
                        groupNames={groupNames ?? []}
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

