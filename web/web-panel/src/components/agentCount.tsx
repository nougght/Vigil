
import styles from "./agentCount.module.css"
import cardSt from "./card.module.css"

export const AgentCount = ({ count, online }: { count: number; online: number }) => {
    return (
        <div className={`${cardSt.card} ${styles.agentsCount}`}>
            <h3>Подключено устройств </h3>
            <div className={styles.valuesContainer}>
                <div className={styles.valueGroup}>
                    <span className={styles.onlineCountValue}>{online}</span>
                    <span className={styles.countLabel}>онлайн</span>
                </div>
                <div className={styles.valueGroup}>
                    <span className={styles.totalCountValue}>/ {count}</span>
                    <span className={styles.countLabel}>всего</span>
                </div>
            </div>
        </div>
    )
}