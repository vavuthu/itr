# ITR

ITR is a test launcher which can launch the test cases in parallel with user controlled fashion having an option for retrying the failed test cases.
ITR should not be tied up with ocs-ci and it can be used with any testing framework or any custom framework having capability to run in a container.

## Why ITR ?

1. ### Execution time:

With the years passing on, our product grows which adds new features for every release. As the product grows, ocs-ci grows for every release which adds new test cases from new features, closed loop etc. As ocs-ci grows, time taken to execute the suites also increases.

Time taken by tier1 2 years ago is not the same as time taken now.

Example:

OCS4-15-Downstream-OCP4-15-VSPHERE6-UPI-FIPS-1AZ-RHCOS-VSAN-3M-3W-tier1 (BUILD ID: 4.15.0-102 RUN ID: 1704364921) —> 19.72 hours

Where as same suite which was run 2 years back
OCS4-9-Downstream-OCP4-9-VSPHERE6-UPI-FIPS-1AZ-RHCOS-VSAN-3M-3W-tier1 (BUILD ID: 4.9.1-254.ci RUN ID: 1639751516) —> 8.58 hours

With the current way of approach, we don’t have any control over time execution.



2. ### Realtime status:

Currently we don’t have any real time status once we start the execution.


## Objective: 

The main objective of ITR is 

1. Reduce the execution time ( suite ): Reduce the execution time of suites by running test cases in parallel. 

2. Real time execution status: should have real time execution status like how many test cases passed so far, how many test cases failed and how many test cases still to run

3. Retry if any test case failed: should have an option to retry if any test cases failed


## How to run:

Since we want to run test cases parallel on a single cluster, negative test cases like osd down, node down/reboot, or any other service operations which affect the cluster will have a direct impact on other test cases which are running parallel. So we need to divide the test cases into 

Positive test cases: Any test cases which DON’T IMPACT’s other test cases which are running parallel.
	
eg: tests/manage/mcg/test_bucket_creation.py::TestBucketCreation::test_bucket_creation[3-CLI-DEFAULT-BACKINGSTORE]


Negative test cases: Any test cases which IMPACT’s other test cases which are running parallel. 

eg: tests/functional/workloads/app/amq/test_amq_node_reboot_and_shutdown.py::TestAMQNodeReboot::test_amq_after_rebooting_node[worker]

Once we have targeted test cases, we will spin up containers for each test case and run them parallel for positive test cases and serially for negative test cases.


Help command for ITR:

```console
$ ./bin/itr -h
2024-04-28T21:32:03.169+0530	INFO	runtime/proc.go:6740	log file: itr_28049-28-09_03-04-09.log
Intelligent Test Runner (ITR) is tool that runs the test cases in parallel with user controlled queues.

Usage:
  itr [flags]

Flags:
  -c, --config_dir string           path to external configuration files that are passed to test framework
  -m, --email string                email to send reports
  -e, --execution string            how to execute the test cases
  -h, --help                        help for itr
  -i, --image string                image name of test framework that should exist in system
  -j, --junit-xml                   Generate JUnit XML report
  -n, --negative_testcases string   Path to negative test cases to run
  -p, --positive_testcases string   Path to positive test cases to run
  -q, --queue-length int            Queue length, number of test cases to run parallelly (default 5)
  -r, --retry int                   number of times to retry the failed test cases
  -s, --subject string              email subject
  -t, --toggle                      Help message for toggle
$ 
```

Run test cases with queuelength of 8 ( 8 test cases runs in parallel ) with retry count as 1 ( failed test case to retry)

```console
./bin/itr -p /home/vijay/VJ/projects/ocs-ci/acceptance_tc_list/acceptance_positive_test_cases -e /home/vijay/VJ/projects/ocs-ci/acceptance_tc_list/how_to_execute_testcase -i localhost/ocsci-testimage1 -c /home/vijay/VJ/clusterdirs/vavuthut1 -q 8 -r 1
```

where as /home/vijay/VJ/projects/ocs-ci/acceptance_tc_list/acceptance_positive_test_cases contains the test cases to run parallelly

```console
$ cat /home/vijay/VJ/projects/ocs-ci/acceptance_tc_list/acceptance_positive_test_cases
tests/functional/pv/pv_services/test_pvc_delete_verify_size_is_returned_to_backendpool.py::TestPVCDeleteAndVerifySizeIsReturnedToBackendPool::test_pvc_delete_and_verify_size_is_returned_to_backend_pool
tests/functional/pv/pv_services/test_dynamic_pvc_accessmodes_with_reclaim_policies.py::TestDynamicPvc::test_rwo_dynamic_pvc[CephFileSystem-Retain]
tests/functional/object/mcg/test_bucket_creation_deletion.py::TestBucketCreationAndDeletion::test_bucket_creation_deletion[3-CLI-DEFAULT-BACKINGSTORE]
tests/functional/object/mcg/test_bucket_creation_deletion.py::TestBucketCreationAndDeletion::test_bucket_creation_deletion[3-OC-DEFAULT-BACKINGSTORE]
tests/functional/object/mcg/test_bucket_creation_deletion.py::TestBucketCreationAndDeletion::test_bucket_creation_deletion[3-S3-DEFAULT-BACKINGSTORE]
tests/functional/object/mcg/test_write_to_bucket.py::TestBucketIO::test_write_file_to_bucket[DEFAULT-BACKINGSTORE]
tests/functional/pv/pvc_snapshot/test_pvc_snapshot.py::TestPvcSnapshot::test_pvc_snapshot[CephFileSystem]
tests/functional/pv/pv_services/test_dynamic_pvc_accessmodes_with_reclaim_policies.py::TestDynamicPvc::test_rwo_dynamic_pvc[CephFileSystem-Delete]
tests/functional/pv/pv_services/test_dynamic_pvc_accessmodes_with_reclaim_policies.py::TestDynamicPvc::test_rwx_dynamic_pvc[CephFileSystem-Retain]
tests/functional/pv/pv_services/test_pvc_assign_pod_node.py::TestPvcAssignPodNode::test_rwx_pvc_assign_pod_node[CephFileSystem]
tests/functional/pv/pvc_clone/test_pvc_to_pvc_clone.py::TestClone::test_pvc_to_pvc_clone[CephFileSystem-None-ReadWriteOnce]
tests/functional/pv/pvc_clone/test_pvc_to_pvc_clone.py::TestClone::test_pvc_to_pvc_rox_clone[CephFileSystem-ReadWriteMany]
tests/functional/pv/pv_services/test_pvc_assign_pod_node.py::TestPvcAssignPodNode::test_rwo_pvc_assign_pod_node[CephFileSystem]
tests/functional/pv/pv_services/test_dynamic_pvc_accessmodes_with_reclaim_policies.py::TestDynamicPvc::test_rwx_dynamic_pvc[CephFileSystem-Delete]
tests/functional/pv/pv_services/test_dynamic_pvc_accessmodes_with_reclaim_policies.py::TestDynamicPvc::test_rwo_dynamic_pvc[CephBlockPool-Retain]
tests/functional/pv/pv_services/test_dynamic_pvc_accessmodes_with_reclaim_policies.py::TestDynamicPvc::test_rwo_dynamic_pvc[CephBlockPool-Delete]
tests/functional/pv/pvc_resize/test_pvc_expansion.py::TestPvcExpand::test_pvc_expansion
tests/functional/pv/pvc_snapshot/test_pvc_snapshot.py::TestPvcSnapshot::test_pvc_snapshot[CephBlockPool]
tests/functional/pv/pv_services/test_raw_block_pv.py::TestRawBlockPV::test_raw_block_pv[Delete]

```

/home/vijay/VJ/projects/ocs-ci/acceptance_tc_list/how_to_execute_testcase contains how to run the test case. 

```console
$ cat /home/vijay/VJ/projects/ocs-ci/acceptance_tc_list/how_to_execute_testcase
run-ci <MY_TEST_CASE> --cluster-path /home/vijay/VJ/clusterdirs/vavuthut1 --ocsci-conf /home/vijay/VJ/clusterdirs/vavuthut1/vSphere7-DC-ECO_VC1 --cluster-name vavuthut1 --disable-environment-checker --resource-checker --kubeconfig /opt/kubeconfig
```

```
$ ls -lrt /home/vijay/VJ/clusterdirs/vavuthut1
total 72096
drwxr-xr-x. 2 vijay vijay       50 Apr 25 15:38 auth
drwxrwxr-x. 2 vijay vijay      116 Apr 25 15:38 terraform_data
-rw-r--r--. 1 vijay vijay      311 Apr 25 15:38 metadata.json
-rw-r--r--. 1 vijay vijay   306409 Apr 25 15:38 bootstrap.ign
-rw-r--r--. 1 vijay vijay     1729 Apr 25 15:38 master.ign
-rw-r--r--. 1 vijay vijay     1729 Apr 25 15:38 worker.ign
-rw-rw-r--. 1 vijay vijay 14025620 Apr 25 15:39 terraform.log
-rw-rw-r--. 1 vijay vijay     1521 Apr 25 17:06 vSphere7-DC-ECO_VC1

$
