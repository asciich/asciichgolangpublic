# commandexecutoriscsi specifications

## Implementation

- The `EnsureInitiorNameConfig` function ensures:
    - The `/etc/iscsi/initiatorname.iscsi` exists.
    - If it does not exsit create it with a proper, new generated InitiatorName provided by the `iscsigeneric` package.
        - Example:
            ```
            # cat /etc/iscsi/initiatorname.iscsi
            InitiatorName=iqn.2016-04.com.open-iscsi:fe84975ad66d
            ```

## Testing

- Use a temporary docker container to test the implementation. This functionality is available from the `dockerutils` package.
