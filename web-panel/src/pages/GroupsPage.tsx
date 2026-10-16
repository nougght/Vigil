// import { Link, useNavigate } from "react-router-dom"
import { useMemo, useState } from "react";
// import { GroupCard } from "../components/groupCard";
import styles from "./groupsPage.module.css"
// import TabBar from "../../../shared/ui/src/components/TabBar.tsx";
import { Checkbox, Dialog } from "radix-ui";
import { CheckIcon } from "@radix-ui/react-icons";
import { useGroups } from "../hooks/useGetGroups.ts";
import { GroupCard } from "../components/groupCard.tsx";
import { useCreateGroup } from "../hooks/useCreateGroup.ts";
import { useAgents } from "../hooks/useGetAgents.ts";
import type { Agent } from "../domain/agent";
import { HideOptions, TextInput } from "@monitoring-system/shared/src/components/TextInput"

interface GroupInfo {
    Name: string
    Description?: string
    Agents?: string[]
}

const AgentsPicker = ({ agents, checked, onChange }: {
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
                                <CheckIcon />
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

export const GroupsPage = () => {
    const [groupInfo, setGroupInfo] = useState<GroupInfo>({
        Name: ""
    })
    const [agentsChecked, setAgentsChecked] = useState<Set<string>>(new Set())
    const [nameWarning, setNameWarning] = useState<string>()

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


    const handleCreate = () => {
        if (!groupInfo.Name) {
            setNameWarning("Название группы не может быть пустым")
            return
        }
        createGroup(
            {
                name: groupInfo.Name,
                description: groupInfo.Description,
                agentIDs: Array.from(agentsChecked)
            },
        )
    }


    const sidePanel = (
        <Dialog.Root>
            <Dialog.Trigger className={styles.createButton}>Новая группа</Dialog.Trigger>
            <Dialog.Portal>
                <Dialog.Overlay className={styles.overlay} />
                <Dialog.Content className={styles.content} asChild>
                    <div>
                        <div>
                            <label className={styles.formLabel} htmlFor="name">Название</label>
                            <TextInput value={groupInfo?.Name}
                                placeholder="Введите название группы"
                                warning={nameWarning}
                                onChange={value => {
                                    value.length > 0 && setNameWarning(undefined);
                                    setGroupInfo({
                                        ...groupInfo,
                                        Name: value
                                    })}
                                 } />
                            <br />
                            <label className={styles.formLabel} htmlFor="description">Описание</label>
                            <TextInput
                                value={groupInfo?.Description}
                                placeholder="Введите описание группы"
                                onChange={v =>
                                    setGroupInfo({
                                        ...groupInfo,
                                        Description: v
                                    })} />

                            {nameWarning && <p className={styles.warning}>{nameWarning}</p>}
                            <AgentsPicker agents={agents ?? []}
                                checked={agentsChecked}
                                onChange={setAgentsChecked} />

                            {_isAgentsPending && <a>Загрузка...</a>}
                            <button onClick={handleCreate}>Создать</button>
                            {isCreatePending && <a>Создание...</a>}
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
                {/*<TabBar tabs={tabs.map((tab) => tab.text)} onSwitch={setActiveTab} activeTab={activeTab} />*/}
                {isGroupsPending && <a>Загрузка</a>}
                <div className={styles.topNav}>
                    <div></div>
                    {sidePanel}
                </div>
                <div className={styles.groupCardsContainer}>{
                    groups && groups.map((g) =>
                        <GroupCard key={g.id} group={g} onClick={() => {
                        }} />
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

