import {ToggleGroup} from "radix-ui";
import styles from "./ToggleChips.module.css"


interface Item<T extends string> {
  value: T
  label: string
}

type BaseProps<T extends string> = {
  items: Item<T>[]
  className?: string;
};

type SingleProps<T extends string> = BaseProps<T> & {
  multiple?: false;
  defaultValue: T;
  // checkedItems: T | null;
  onChange: (checked: T ) => void;
};

type MultiProps<T extends string> = BaseProps<T> & {
  multiple: true;
  defaultValue: T[];
  // checkedItems: T[];
  onChange: (checked: T[]) => void;
};

type Props<T extends string> = MultiProps<T> | SingleProps<T>;


export function ToggleChips<T extends string>(
  {
    items,
    className,
    defaultValue,
    multiple,
    onChange,
  }: Props<T>) {

  const item = <>
    {
      items.map(item =>
        <ToggleGroup.Item
          className={styles.ToggleChip}
          value={item.value}
        >
          <span>{item.label}</span>
        </ToggleGroup.Item>
      )
    }
  </>

  return (
    <>
      {
        multiple ?
          <ToggleGroup.Root
            className={styles.ToggleGroup + ` ${className}`}
            type="multiple"
            aria-label="Text alignment"
            onValueChange={(values) => onChange(values as T[])}
          >
            {
              item
            }
          </ToggleGroup.Root>
          :
          <ToggleGroup.Root
            className={styles.ToggleGroup + ` ${className}`}
            type="single"
            aria-label="Text alignment"
            defaultValue={defaultValue}
            onValueChange={(value) => onChange(value as T)}
          >
            {
              item
            }
          </ToggleGroup.Root>
      }
    </>
  )
}