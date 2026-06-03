# Journey Flow

| Field         | Value                                   |
| ------------- | --------------------------------------- |
| Journey       | JRN.<domain>.<goal>.<channel>.<variant> |
| Platform      | TBD(owner)                              |
| Auth state    | TBD(owner)                              |
| Locale/device | TBD(owner)                              |

## Flow

```mermaid
flowchart TD
    A[Entry] --> B[Step]
    B --> C{Decision}
    C -->|Happy path| D[Success]
    C -->|Recovery| E[Recovery]
    E --> B
```

## Nodes

| Node | Screen/touchpoint | User action | System response | State/data |
| ---- | ----------------- | ----------- | --------------- | ---------- |
| A    | TBD(owner)        | TBD(owner)  | TBD(owner)      | TBD(owner) |

## Decisions

| Decision   | Branches   | Rule       | Recovery   |
| ---------- | ---------- | ---------- | ---------- |
| TBD(owner) | TBD(owner) | TBD(owner) | TBD(owner) |

## Accessibility

| Area          | Requirement | Check      |
| ------------- | ----------- | ---------- |
| Keyboard      | TBD(owner)  | TBD(owner) |
| Screen reader | TBD(owner)  | TBD(owner) |
| Error state   | TBD(owner)  | TBD(owner) |
