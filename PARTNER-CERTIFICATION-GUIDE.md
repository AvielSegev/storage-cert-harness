# OpenShift Virtualization CSI Partner Certification Guide
## Level 3: Performance Validated

**Red Hat Ecosystem Certification**  
**Version:** 1.0  
**Last Updated:** September 2026

---

## Overview

This guide provides instructions for running the OpenShift Virtualization CSI Partner Certification test suite to achieve **Level 3 (Performance Validated)** certification status. The certification test harness validates your storage solution against 11 critical performance and reliability test requirements.

## Prerequisites

Before beginning certification testing, ensure you have:

1. **OpenShift Cluster** (version 4.20 or later) with:
   - OpenShift Virtualization installed and configured
   - At least 3 worker nodes for scale testing
   
2. **Storage Backend** configured with your storage class under test

3. **Container Runtime** with access to Red Hat Quay (podman or docker)

4. **Certification Materials** from Red Hat:
   - `catalog.json` — test specifications (available in this repository)
   - `thresholds.json` — performance thresholds (**confidential**, provided via secure email)

## Required Files

### 1. Test Catalog (Public)

Download the latest catalog from this repository:

```bash
curl -O https://raw.githubusercontent.com/openshift-storage-cert/storage-cert-harness/main/catalog.json
```

### 2. Performance Thresholds (Confidential)

The `thresholds.json` file containing SLA gates will be provided to you separately via encrypted email by your Red Hat Ecosystem Certification contact. **Do not share this file publicly.**

### 3. Test Plans (Included)

Two customer test plans are required for Level 3 certification:
- `level3-customer-load-80.yaml` — Data path test (storage I/O latency under load)
- `level3-customer-no-load-80.yaml` — Control path, scale, and lifecycle tests

Download both plans:

```bash
curl -O https://raw.githubusercontent.com/openshift-storage-cert/storage-cert-harness/main/plans/level3-customer-load-80.yaml
curl -O https://raw.githubusercontent.com/openshift-storage-cert/storage-cert-harness/main/plans/level3-customer-no-load-80.yaml
```

### 4. Storage Backend Configuration (Required)

You must create a `backends.yaml` file defining your storage configuration. See the **Storage Backend Configuration** section below for detailed instructions and examples.

## Storage Backend Configuration

Before running the certification tests, you must define your storage backend configuration. Create a `backends.yaml` file that specifies your storage class and vendor details:

```yaml
schema_version: "1"
backends:
  - name: my-storage-backend
    vendor: <your-vendor>  # e.g., netapp-trident, pure-storage, dell-emc
    storage_class: <your-storage-class-name>
    snapshot_class: <your-snapshot-class-name>  # Optional
    params: {}  # Vendor-specific parameters (if needed)
    secrets: {}  # Authentication credentials (if needed)
```

**Example for a basic CSI driver:**

```yaml
schema_version: "1"
backends:
  - name: partner-storage
    storage_class: fast-ssd-csi
    snapshot_class: fast-ssd-snapclass
```

**Important:** Replace `<your-storage-class-name>` with the actual StorageClass name deployed in your OpenShift cluster. The backend name (`my-storage-backend` or `partner-storage` in the examples above) will be used in the `--backend` flag when running tests.

For backends requiring authentication, see `backends.example.yaml` in this repository for credential reference examples.

## Running the Certification Tests

### Step 1: Obtain the Harness Container Image

Pull the latest certification harness image from Red Hat Quay:

```bash
podman pull quay.io/virtarraycert/storage-cert-harness:latest
```

Or specify a specific version:

```bash
podman pull quay.io/eco-special-projects/storage-cert-harness:<Version>
```

### Step 2: Configure Test Plans

Edit both test plan files to specify your cluster details:

**`level3-customer-no-load-80.yaml`** requires the some placeholders to be replaced, e.g:

```yaml
TR-VIRT-007:
  node: "<TARGET_NODE_FOR_FAILURE_RECOVERY>"  # Replace with actual worker node name

TR-VIRT-008:
  target_node: "<TARGET_NODE_TO_DRAIN>"       # Replace with actual worker node name

TR-VIRT-013:
  node_name: "<TARGET_NODE_FOR_SINGLE_NODE_CEILING>"  # Replace with actual worker node name
```

**`level3-customer-load-80.yaml`** typically requires no changes, but you may optionally adjust:

```yaml
TR-STOR-002:
  # num_vms: <NUMBER_OF_TEST_VMS>  # Uncomment and set if needed
```

### Step 3: Run Test Plan 1 — Data Path (Load 80)

Execute the storage I/O latency test under load:

```bash
podman run --rm \
  -v $(pwd)/catalog.json:/catalog.json:ro \
  -v $(pwd)/thresholds.json:/thresholds.json:ro \
  -v $(pwd)/backends.yaml:/backends.yaml:ro \
  -v $(pwd)/level3-customer-load-80.yaml:/plan.yaml:ro \
  -v $(pwd)/kubeconfig:/kubeconfig:ro \
  -v $(pwd)/results:/results \
  -e KUBECONFIG=/kubeconfig \
  quay.io/eco-special-projects/storage-cert-harness:latest \
  run \
    --catalog /catalog.json \
    --thresholds /thresholds.json \
    --backends /backends.yaml \
    --backend my-storage-backend \
    --plan /plan.yaml \
    --output /results/load-80
```

**Note:** Replace `my-storage-backend` with the backend name you defined in your `backends.yaml` file.

### Step 4: Run Test Plan 2 — Control Path, Scale & Lifecycle

Execute the remaining 10 test requirements:

```bash
podman run --rm \
  -v $(pwd)/catalog.json:/catalog.json:ro \
  -v $(pwd)/thresholds.json:/thresholds.json:ro \
  -v $(pwd)/backends.yaml:/backends.yaml:ro \
  -v $(pwd)/level3-customer-no-load-80.yaml:/plan.yaml:ro \
  -v $(pwd)/kubeconfig:/kubeconfig:ro \
  -v $(pwd)/results:/results \
  -e KUBECONFIG=/kubeconfig \
  quay.io/eco-special-projects/storage-cert-harness:latest \
  run \
    --catalog /catalog.json \
    --thresholds /thresholds.json \
    --backends /backends.yaml \
    --backend my-storage-backend \
    --plan /plan.yaml \
    --output /results/no-load-80
```

**Note:** Replace `my-storage-backend` with the backend name you defined in your `backends.yaml` file.

## Understanding Results

### Exit Codes

- **0** — All tests passed certification criteria
- **1** — One or more tests failed to meet SLA thresholds
- **2** — Harness error (invalid configuration, cluster access issues)

### Output Artifacts

Each test run produces the following artifacts in the `--output` directory:

```
results/
├── load-80/
│   ├── report.json          # Machine-readable test results
│   ├── report.md            # Human-readable certification report
│   └── artifacts/           # Raw test tool outputs
└── no-load-80/
    ├── report.json
    ├── report.md
    └── artifacts/
```

### Certification Criteria

To achieve **Level 3 (Performance Validated)** certification:

1. **Both test plans must complete successfully** (exit code 0)
2. **All 11 test requirements must pass** their SLA thresholds
3. **Test reports must be submitted** to Red Hat for verification

## Test Requirements Summary

Level 3 certification validates the following capabilities:

| Test ID | Category | Test Name |
|---------|----------|-----------|
| TR-STOR-002 | Data path | Storage I/O Latency (under load) |
| TR-STOR-006 | Control path | Provisioning Density (1000 PVCs) |
| TR-VIRT-001 | Scale | Boot Storm (100-1000 VMs) |
| TR-VIRT-004 | Scale | Parallel Live Migration (50 VMs) |
| TR-VIRT-007 | Life-cycle | VM HA Restart (node failure recovery) |
| TR-VIRT-008 | Life-cycle | Node Drain (VM evacuation) |
| TR-VIRT-010 | Control path | Snapshot Creation (batch p99) |
| TR-VIRT-013 | Scale | Per-Node Ceiling (max VMs/volumes per node) |
| TR-VIRT-018 | Control path | Clone at Scale (100 VDI clones) |
| TR-VIRT-019 | Life-cycle | Parallel VM Lifecycle (20 VMs) |
| TR-VIRT-027 | Control path | Snapshot Depth (250 snapshots per volume) |

