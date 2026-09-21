import { useEffect, useState } from "react";
import type { ActivityUpdate } from "../domain/activity";

import styles from "./agentActivity.module.css"



export const AgentActivity = ({ activity, agentID }: { activity?: ActivityUpdate, agentID: string }) => {
    const [isImgMax, setIsImgMax] = useState(false)
    const [width, setWidth] = useState("30%")

    useEffect(() => {
        setWidth(isImgMax ? '90%' : `40%`)
    }, [isImgMax])

    const handleResize =() => {
        console.log("resize stream - ", isImgMax)
        setIsImgMax(!isImgMax)
    }

    return (
        <div className={styles.agentActivity}
            style={{
                '--width': width,
            } as React.CSSProperties}>

            <h3>Активное приложение</h3>
            <p>{activity?.title != null ? activity?.title : "Нет данных"}</p>
            <h3>Трансляция экрана</h3><br />
            <img src={`http://monitoring.nought.ru/api/v1/agents/${agentID}/frames`} onClick={handleResize}/>
        </div>
    )
}