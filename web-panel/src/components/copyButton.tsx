import { CopyToClipboard} from "@monitoring-system/shared"
import copyIcon from "../assets/copy.svg"
import styles from "./copyButton.module.css"

export const CopyButton = (props: { text: string }) => {
    const handleClick = () => {
        CopyToClipboard(props.text)
    }
    return (
        <button className={styles.copyButton} onClick={handleClick}>
            <img src={copyIcon} width="16" height="16"/>
        </button>
    )
}