import {Eye, EyeOff} from 'lucide-react';

import styles from './TextInput.module.css'
import {useEffect, useRef, useState} from 'react';

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
  multiline?: boolean
  className?: string
}

export const TextInput = (props: TextInputProps) => {
  const [hide, setHide] = useState(props.hide == HideOptions.AlwaysHide || props.hide == HideOptions.ToggleHide)
  // const [height, setHeight] = useState(40)
  const maxHeight = 250;
  // const areaHeight = `${height}px`

  const areaRef = useRef<HTMLTextAreaElement>(null)

  useEffect(
    () => {
      if (areaRef.current == null) {
        return
      }

      const textarea = areaRef.current;
      // 1. Сбрасываем высоту, чтобы корректно рассчитать scrollHeight при удалении текста
      textarea.style.height = 'auto';
      // 2. Устанавливаем высоту под контент
      textarea.style.height = `${textarea.scrollHeight > maxHeight ? maxHeight : textarea.scrollHeight}px`;

    }, [props.value]
  )
  return (
    <>
      <div className={(props.className ? `${props.className} ` : "") + styles.TextInputContainer + (props.warning ? ` ${styles.warning}` : "")}>
        {
          props.multiline ?
            <textarea className={styles.TextInput}
                      ref={areaRef}
                      style={{
                        '--height': "40px",
                        '--max-height': `${maxHeight}px`
                      } as React.CSSProperties}
                      placeholder={props.placeholder}
                      value={props.value}
                      onChange={e =>
                        props.onChange(e.target.value)}
            />
            : <input className={styles.TextInput}
                     type={hide ? "password" : "text"}
                     placeholder={props.placeholder}
                     value={props.value}
                     onChange={e =>
                       props.onChange(e.target.value)}
            />
        }
        {
          props.hide == HideOptions.ToggleHide
          && <HideToggle
                hide={hide}
                onToggle={() => {
                  setHide(!hide)
                }}
            />
        }

      </div>
      {
        props.warning &&
          <p className={styles.warning}>{props.warning}</p>
      }
    </>
  )
}

const HideToggle = ({hide, onToggle}: { hide: boolean, onToggle: () => void }) => {
  return (
    <div className={styles.HideToggle} onClick={onToggle}>
      {hide ? <Eye size={15}/> : <EyeOff size={15}/>}
    </div>
  )
}