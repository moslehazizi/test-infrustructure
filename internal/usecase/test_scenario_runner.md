Test scenario categorization based on how it should be executed:

- A (total test service + total execution time): smoke, load, soak, spike.
- B (step execution + increase rate): stress.
- C: (total test service + step execution + increase rate): scalability.
- D: (total test service + step execution + increase rate + decrease rate): recovery.

## Group A

```mermaid
flowchart TD
    start((START))
    totalCountCond{all test services created?}
    executionTimeCond{execution time is up}
    wait[wait for specified time]
    createTestService[Create total required test services]
    shutdownServices[Shut down all test services]
    END((END))

    start --> totalCountCond
    totalCountCond -->|no| createTestService --> start
    totalCountCond -->|yes| executionTimeCond

    executionTimeCond -->|no| wait --> executionTimeCond
    executionTimeCond -->|yes| shutdownServices

    shutdownServices --> END
```
