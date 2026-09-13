
import styles from "./agentCount.module.css"
import cardSt from "./card.module.css"

export const AgentCount = ({ count, online }: { count: number; online: number }) => {
    return (
        <div className={`${cardSt.card} ${styles.agentsCount}`}>
            <h3>Подключено устройств </h3>
            <div>
                <div className={styles.valuesContainer}>
                    <span className={styles.onlineCountValue}>{online}</span>
                    <span className={styles.totalCountValue}>/ {count}</span>
                </div>
                <div className={styles.countLabelsContainer}>
                    <span className="online-label">онлайн</span>
                    <span className="total-label">всего</span>
                </div>
            </div>
        </div>
    )
}