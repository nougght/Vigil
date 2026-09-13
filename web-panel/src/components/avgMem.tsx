
import styles from "./avg.module.css"
import cardSt from "./card.module.css"

export const AvgMem = ({ value }: { value: number }) => {
    return (
        <div className={`${cardSt.card} ${styles.avg}`}>
            <h3>Средняя нагрузка ОЗУ</h3>
            <div>
                <span>{value.toFixed(2)}%</span>
            </div>
        </div>
    )
}