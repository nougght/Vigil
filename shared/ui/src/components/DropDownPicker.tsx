import {useState} from "react";
import {DropdownMenu} from "radix-ui";
import styles from "./dropDownPicker.module.css"
import {CheckIcon, ChevronDownIcon, DotFilledIcon} from "@radix-ui/react-icons";
import {TextInput} from "./TextInput";

interface Item<T extends string> {
  value: T
  label: string
}

type BaseProps<T extends string> = {
  items: Item<T>[]
  triggerPlaceholder: string
  needSearch?: boolean;
  className?: string;
};

type SingleProps<T extends string> = BaseProps<T> & {
  multiple?: false;
  checked: T | null;
  onChange: (checked: T | null) => void;
};

type MultiProps<T extends string> = BaseProps<T> & {
  multiple: true;
  checked: T[];
  onChange: (checked: T[]) => void;
};


type Props<T extends string> = SingleProps<T> | MultiProps<T>

export function DropDownPicker<T extends string>(
  {
    items,
    checked,
    triggerPlaceholder,
    onChange,
    multiple,
    needSearch,
    className
  }: Props<T>) {
  const [searchItem, setSearchItem] = useState<string>()

  const toggle = (value: T, itemChecked?: boolean) => {
    if (multiple === true) {
      onChange(itemChecked == true ?
        [...checked, value] :
        checked.filter(item => item !== value)
      )
    } else {
      onChange(value)
    }
  }

  const label =
    multiple === true ?
      checked?.length === 0
        ? triggerPlaceholder
        : checked.length == 1
          ? items.filter((o) => checked.includes(o.value)).map((o) => o.label)[0]
          : `Выбрано: ${checked.length}`
      :
      checked == null ?
        triggerPlaceholder :
        items.filter((o) => o.value === checked).map((o) => o.value)[0]

  const visibleItems = needSearch && searchItem ?
    items.filter(item => item.label.includes(searchItem)) :
    items

  return (
    <DropdownMenu.Root>
      <DropdownMenu.Trigger className={className + ` ${styles.DropDownTrigger}`} disabled={false}>
        <span>{label}</span>
        <ChevronDownIcon aria-hidden/>
      </DropdownMenu.Trigger>

      <DropdownMenu.Portal>
        <DropdownMenu.Content className={styles.DropDownContent}
                              align="start"
                              sideOffset={4}
                              // alignOffset={100}
                              side="right"
                              loop>
          <TextInput
            className={styles.searchField}
            onChange={(value) => {
              setSearchItem(value);
            }}
            value={searchItem}
            placeholder={"Поиск"}

          />
          {multiple ? (
            visibleItems.map((item) => (
              <DropdownMenu.CheckboxItem
                key={item.value}
                className={styles.CheckboxItem}
                checked={checked.includes(item.value)}
                onCheckedChange={(checked) =>
                  toggle(item.value, checked)
                }
                onSelect={(e) => e.preventDefault()}
              >
                <span className={styles.CheckboxIndicator}>
                  <DropdownMenu.ItemIndicator>
                    <CheckIcon height={16} width={16}/>
                  </DropdownMenu.ItemIndicator>
                </span>
                {item.label}
              </DropdownMenu.CheckboxItem>
            ))
          ) : (
            <DropdownMenu.RadioGroup
              value={checked ?? ""}
              onValueChange={(v) => toggle(v as T)}
            >
              {visibleItems.map((item) => (
                <DropdownMenu.RadioItem
                  key={item.value}
                  value={item.value}
                  className={styles.RadioItem}
                  onSelect={(e) => e.preventDefault()}
                >
                <span className={styles.RadioIndicator}>
                  <DropdownMenu.ItemIndicator>
                    <DotFilledIcon height={16} width={16}/>
                  </DropdownMenu.ItemIndicator>
                </span>
                  {item.label}
                </DropdownMenu.RadioItem>
              ))}
            </DropdownMenu.RadioGroup>
          )}
        </DropdownMenu.Content>
      </DropdownMenu.Portal>
    </DropdownMenu.Root>
  );
}