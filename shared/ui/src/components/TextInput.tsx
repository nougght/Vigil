
import { Eye, EyeOff } from 'lucide-react';

import styles from './TextInput.module.css'
import { useState } from 'react';

export const HideOptions = {
    NotHide: "not_hide",
    AlwaysHide: "hide",
    ToggleHide: "toggle_hide"
} as const;

export type HideOptions = (typeof HideOptions)[keyof typeof HideOptions];


interface TextInputProps {
    value?: string
    placeholder?: string
    onChange: (value: string) => void
    warning?: string
    hide?: HideOptions
}

export const TextInput = (props: TextInputProps) => {
    const [hide, setHide] = useState(props.hide == HideOptions.AlwaysHide || props.hide == HideOptions.ToggleHide)

    return (
        <>
        <div className={styles.TextInputContainer + (props.warning ? ` ${styles.warning}` : "")}>
            <input className={styles.TextInput}
                type={hide ? "password" : "text"}
                placeholder={props.placeholder}
                value={props.value}
                onChange={e =>
                    props.onChange(e.target.value)}
            />
            {
                props.hide == HideOptions.ToggleHide
                && <HideToggle
                    hide={hide}
                    onToggle={() => {
                        setHide(!hide)
                    }}
                />
            }

        </div >
        {
            props.warning && 
            <p className={styles.warning}>{props.warning}</p>
        }
        </>
    )
}

const HideToggle = ({ hide, onToggle }: { hide: boolean, onToggle: () => void }) => {
    return (
        <div className={styles.HideToggle} onClick={onToggle}>
            {hide ? <Eye size={15} /> : <EyeOff size={15} />}
        </div>
    )
}