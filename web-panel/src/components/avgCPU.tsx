
import styles from "./avg.module.css"
import cardSt from "./card.module.css"

export const AvgCPU = ({ value }: { value: number }) => {
    return (
        <div className={`${cardSt.card} ${styles.avg}`}>
            <h3>Средняя нагрузка CPU</h3>
            <div>
                <span>{value.toFixed(2)}%</span>
            </div>
        </div>
    )
}