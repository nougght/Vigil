import type {ReactNode} from "react";
import styles from "./Table.module.css"
import {ArrowDown, ArrowUp} from "lucide-react";
import {CheckIcon} from "@radix-ui/react-icons";
import {Checkbox} from "radix-ui";


export interface TableColumn<T> {
  key: string;      // unique key
  header: string;
  sortable?: boolean;
  fieldRenderFunc?: (row: T) => ReactNode;  // custom render for column fields
}

interface TableProps<T> {
  rows: T[];
  columns: TableColumn<T>[];
  selected: Set<string>;
  sortColumn?: string;             // column key for sort
  sortOrder?: 'asc' | 'desc';
  onSelect?: (ids: string[], all?: boolean) => void
  onSort?: (columnKey: string) => void;
  onRowClick?: (id: string) => void;
}

const TableCheckbox = ({id, checked, toggleChecked}: { id: string, checked: boolean, toggleChecked: () => void }) => {
  return (
    // TODO: move to shared components
    <Checkbox.Root id={id}
                   className={styles.checkboxRoot}
                   checked={checked}
                   onCheckedChange={() => toggleChecked()}
    >
      <Checkbox.Indicator className={styles.checkboxInd}>
        <CheckIcon/>
      </Checkbox.Indicator>
    </Checkbox.Root>
  )
}

export function Table<T extends { id: string }>(props: TableProps<T>) {
  return (
    <div className={styles.TableWrapper}>
      <table className={styles.Table}>
        <thead>

        <tr className={styles.HeaderRow}>
          <th>
            <TableCheckbox id="all"
                           checked={props.selected.size == props.rows.length}
                           toggleChecked={() => props.onSelect?.([], true)}/>
          </th>
          {
            props.columns.map((column, key) => {
              const isSortActive = props.sortColumn == column.key;
              return (
                <th key={key}
                    className={styles.Header}
                    onClick={() => {
                      props.onSort?.(column.key);
                      console.log(isSortActive);
                      console.log(props.sortOrder);
                    }
                    }
                >
                  {column.header}
                  {isSortActive &&
                    (
                      props.sortOrder === 'desc' && <ArrowDown size={16}/> ||
                      props.sortOrder === 'asc' && <ArrowUp size={16}/>
                    )
                  }
                </th>
              )
            })
          }
        </tr>
        </thead>
        <tbody>
        {
          props.rows.map((row: T) => {
            return (
              <tr onClick={() => {props.onRowClick?.(row.id)}}>
                <td onClick={e => e.stopPropagation()}><TableCheckbox id={row.id}
                                   checked={props.selected.has(row.id)}
                                   toggleChecked={() => props.onSelect?.([row.id])}/>
                </td>
                {
                  props.columns.map((column, i) => {
                    const value = row[column.key as keyof T]
                    return (
                      <td key={i}>
                        {
                          column.fieldRenderFunc?.(row) ??
                          (value == null ? "" : String(value))
                        }
                      </td>
                    )
                  })
                }
              </tr>
            )
          })
        }
        </tbody>
      </table>
    </div>
  )
}