# workload-sleeper
This is aimed towards creating a controller that will make workloads sleep

### Example

```yaml
apiVersion: greenworkload.platform.io/v1beta1
kind: SleepSchedule
metadata:
  labels:
    app.kubernetes.io/name: workload-sleeper
    app.kubernetes.io/managed-by: kustomize
  name: sleepschedule-sample
spec:
  schedule:
    pauseScheduled: false
    sleepEndDate: "2024-08-04"
    sleepStartDate: "2024-08-03"
    sleepStartTime: "20:00"
    sleepEndTime: "20:50"
    timeZone: "UTC"

```