# DevOps Mini Project — Go API บน Kubernetes ที่ดูแลเอง

> โปรเจค DevOps แบบ end-to-end: นำ Go REST API มา containerize แล้ว deploy ขึ้น
> **k3s** Kubernetes cluster ที่ตั้งและดูแลเอง (provision ด้วย **Ansible**) รองรับหลาย environment พร้อม
> **CI/CD pipeline**, **security scanning**, **HTTPS อัตโนมัติ**, **monitoring**, **autoscaling**
> และ **image promotion** ที่ build ครั้งเดียวแล้วใช้ image ตัวเดิมจาก **dev → production**

![CI](https://github.com/thanutgit/DevOpsMiniProject/actions/workflows/ci.yaml/badge.svg)

🔗 **Live:** https://45.150.128.180.nip.io

> ใช้ Aiven MySQL free tier ซึ่งจะปิดตัวเองเมื่อไม่มีการใช้งาน หากเข้าไม่ได้แปลว่า
> database อยู่ในสถานะพัก

---

## ภาพรวม (Overview)

โปรเจคนี้สาธิต workflow การ deploy ที่ใกล้เคียงงานจริง สำหรับ web service ขนาดเล็ก —
ตั้งแต่ source code จนถึงแอปที่รันอยู่จริงผ่าน HTTPS มี monitoring และ auto-scale บน
infrastructure จริง (ไม่ใช่แค่ demo บนเครื่อง local)

Container image ตัวเดียวกันถูกใช้รัน 2 บทบาท โดยเลือกตอน runtime ผ่าน environment
variable:

- **`server`** — เปิดให้บริการ HTTP API
- **`migrator`** — รัน database migration ครั้งเดียวแล้วจบ (Kubernetes `Job`)

แยกออกเป็น 2 environment ที่อิสระจากกัน:

| Environment | Infrastructure | Database | เข้าถึง |
|---|---|---|---|
| **dev** | k3s บน VirtualBox VM (local) | Aiven MySQL (database `*_dev`) | HTTP ผ่าน IP ของ VM |
| **prd** | k3s บน VPS 2 เครื่อง (control-plane 1 + worker 1) | Aiven MySQL (database production) | HTTPS ผ่าน `45.150.128.180.nip.io` |

> dev ไม่มี TLS เพราะอยู่ใน IP ภายใน Let's Encrypt ยิงเข้ามาตรวจ (HTTP-01) ไม่ถึง

---

## สถาปัตยกรรม (Architecture)

### Runtime architecture (production)

```mermaid
flowchart TD
    User([User / Browser]) -->|HTTPS :443| Traefik[Traefik Ingress]
    User -.->|HTTP :80 ถูก redirect 308| Traefik
    LE[Let's Encrypt] -.HTTP-01 challenge.-> Traefik
    CM[cert-manager] -.ขอ/ต่ออายุ TLS cert.-> LE
    CM -.TLS secret.-> Traefik
    Traefik -->|path /| Svc[app-service ClusterIP]
    Traefik -->|path /grafana| Graf[Grafana]
    Svc --> P1[app-api pod]
    Svc --> P2[app-api pod]
    Svc --> P3[app-api pod]
    HPA[HorizontalPodAutoscaler<br/>CPU 85%, 3-6 replicas] -.scales.-> Svc
    P1 --> DB[(Aiven MySQL)]
    P2 --> DB
    P3 --> DB
    Job[migrate Job<br/>รันก่อน app] --> DB

    subgraph Cluster["k3s cluster (master + worker)"]
        Traefik
        CM
        Svc
        P1
        P2
        P3
        HPA
        Job
        Graf
    end
```

### CI/CD pipeline

```mermaid
flowchart LR
    Push[Push เข้า dev branch] --> CI{CI<br/>test + security scan}
    CI -->|ผ่าน| Build[Build image ครั้งเดียว<br/>devopsminiproject-dev:TAG]
    Build --> DevDeploy[Deploy dev + ทดสอบ]
    DevDeploy --> PR[Pull Request dev to main<br/>CI ต้องผ่านก่อน merge]
    PR -->|merge| Promote[Promote จาก main<br/>copy image ไป -prd:TAG<br/>ตรวจ digest ตรงกัน]
    Promote --> PrdDeploy[Deploy prd ด้วย TAG เดิม]
    PrdDeploy --> K8s[apply manifests<br/>รัน migrate Job แล้ว rollout]
```

---

## Tech Stack

| ส่วน | เครื่องมือ |
|---|---|
| **Language / Framework** | Go 1.26, Fiber v3, GORM |
| **Containerization** | Docker (multi-stage build บน alpine), GitHub Container Registry (GHCR) |
| **Orchestration** | Kubernetes (k3s), Traefik Ingress, Helm |
| **Infrastructure as Code** | Ansible (roles: common, k3s_server, k3s_agent, platform), Ansible Vault |
| **TLS** | cert-manager, Let's Encrypt (HTTP-01) |
| **Security** | Trivy (scan repo + image), Kubernetes securityContext (non-root, read-only filesystem) |
| **CI/CD** | GitHub Actions, self-hosted runners, GitHub Environments |
| **Monitoring / Alerting** | Prometheus, Grafana, Alertmanager (kube-prometheus-stack), แจ้งเตือนผ่าน Discord |
| **Testing** | Go testing (table-driven unit tests), k6 (load testing) |
| **Database** | MySQL (Aiven) |
| **Infrastructure** | VPS (production), VirtualBox (development) |

---

## ฟีเจอร์เด่น (Key Features)

- **Provision cluster ด้วย Ansible** — จาก Ubuntu เปล่าเป็น k3s cluster (master +
  worker) ในคำสั่งเดียว แบ่งเป็น 3 role: `common` (hostname, ปิด cloud-init network,
  ลบ search domain, firewall), `k3s_server` และ `k3s_agent` (worker ดึง token จาก
  master อัตโนมัติผ่าน `hostvars`) pin version ของ k3s ให้ตรงกับ prd, ซ่อน token
  ด้วย `no_log` และ **idempotent** (รันซ้ำได้ `changed=0`) ทดสอบบน VM โดยถอน k3s
  ออกแล้วสร้าง cluster ใหม่จาก playbook จากนั้นนำมาจัดการ **prd (VPS) จริง** โดยรัน
  `--check --diff` ก่อนเพื่อหา config drift แล้วค่อย apply ใช้ inventory แยกต่อ
  environment และเปิด/ปิด firewall ได้ด้วยตัวแปร (`manage_ufw`)
- **Platform เป็น code ทั้งหมด** — role `platform` ติดตั้ง helm (pin version),
  kube-prometheus-stack และ cert-manager ด้วย `kubernetes.core.helm` (pin chart version)
  และ apply ClusterIssuer + alert rules ด้วย `kubernetes.core.k8s` โดยอ่าน values และ
  manifest จาก `k8s/platform/` ใน repo โดยตรง (git เป็นแหล่งความจริงที่เดียว) Discord
  webhook เก็บใน **Ansible Vault** แล้วรวมเข้ากับ values **ในหน่วยความจำ** ไม่มีไฟล์ secret
  บนดิสก์ เปิด/ปิด cert-manager ต่อ environment ได้ และ pin collection ไว้ใน
  `requirements.yaml` นำมาครอบ prd ที่เคยติดตั้งด้วยมือได้ด้วย `--check --diff` ที่ผล
  เป็น `ok` ทั้งหมด (ของจริงตรงกับ repo ทุกอย่าง ไม่มี upgrade ซ้ำ)
- **Supply chain ของสคริปต์ติดตั้ง** — สคริปต์ติดตั้ง k3s ถูกตรวจด้วย sha256 ที่ pin
  ไว้คู่กับ `k3s_version` ถ้าสคริปต์ถูกแก้จะไม่ถูกรันด้วย root
- **SSH keys เป็น code** — public key ของทุกเครื่องที่เข้าได้อยู่ใน repo และถูกใส่ด้วย
  `authorized_key` (comment ตามชื่อไฟล์) เพิ่ม/ถอนสิทธิ์ผ่าน PR
- **Server hardening** — firewall (ufw) แบบ deny เป็นค่าเริ่มต้น เปิดสาธารณะแค่
  22/80/443 และ SSH รับเฉพาะ key (ปิด password และ keyboard-interactive, root เข้าได้
  ด้วย key เท่านั้น) ตรวจ syntax ด้วย `sshd -t` ก่อนวางไฟล์ และยืนยันค่าที่ใช้จริงด้วย
  `sshd -T`
- **Image promotion (build once, deploy many)** — build image ครั้งเดียวบน dev แล้ว
  promote image ตัวเดิมไป prd ด้วย `docker buildx imagetools create` โดยไม่ build ใหม่
  workflow ตรวจว่า **digest ของ dev และ prd ตรงกัน** เพื่อยืนยันว่า prd รัน image
  ตัวเดียวกับที่ทดสอบแล้ว, promote ได้จาก `main` เท่านั้น (ต้องผ่าน PR + CI) และ
  ห้ามเขียนทับ tag ที่มีอยู่แล้วบน prd
- **Security scanning ใน CI (Trivy)** — สแกน repository (ช่องโหว่ของ dependency,
  secret ที่หลุดในโค้ด, misconfiguration ของ Dockerfile/Kubernetes manifest) และสแกน
  container image ที่ build จริง พบช่องโหว่ระดับ HIGH/CRITICAL ที่มีแพตช์แล้ว CI จะ fail
  และเป็น required check ก่อน merge เข้า `main`
- **Container hardening** — รันเป็น non-root user (UID 10001), root filesystem แบบ
  read-only, ตัด Linux capabilities ทั้งหมด, ปิด privilege escalation และใช้
  seccomp profile `RuntimeDefault`
- **Multi-stage Docker build** — แยก stage build ออกจาก runtime ทำให้ image เล็กลงมาก
  และฝัง build metadata (version, build time, commit SHA) ลงใน binary ผ่าน `-ldflags`
  แสดงผลที่ endpoint `/about`
- **Database migration เป็น Kubernetes Job** ที่รันจนเสร็จก่อนจะ rollout แอป โดยใช้
  image ตัวเดียวกันในโหมด `migrator`
- **แยก health probes ชัดเจน**:
  - `/livez` (liveness) — เช็คแค่ว่า process ยังทำงาน **ไม่แตะ DB** เพื่อไม่ให้
    DB หลุดชั่วคราวทำให้ pod ถูก restart โดยไม่จำเป็น
  - `/readyz` (readiness) — เช็ค database ด้วย เพื่อส่ง traffic เข้าเฉพาะ pod ที่
    พร้อมให้บริการจริง
- **HTTPS อัตโนมัติ** — cert-manager ขอและต่ออายุ certificate จาก Let's Encrypt
  ทดสอบกับ staging issuer ก่อนเปลี่ยนเป็น production เพื่อเลี่ยง rate limit
  และ redirect HTTP → HTTPS ทุก request ด้วย Traefik Middleware
- **Horizontal Pod Autoscaler** — scale `app-api` จาก 3 ถึง 6 replicas ที่ CPU
  85% (ทดสอบด้วย k6)
- **Unit tests แบบ table-driven** รันใน CI ก่อน merge ทุกครั้ง
- **การแยก config ตาม environment** — secret แยกผ่าน GitHub Environments,
  ConfigMap และ Ingress แยกไฟล์ต่อ environment
- **แยก platform กับ application** — ของที่ตั้งครั้งเดียวต่อ cluster (`k8s/platform/`)
  แยกจาก manifest ที่ deploy ทุก release
- **Monitoring stack** — Prometheus + Grafana ติดตั้งผ่าน Helm เปิดที่ `/grafana`
- **Alert เมื่อ certificate ใกล้หมดอายุ** — ให้ Prometheus เก็บ metrics ของ cert-manager
  ผ่าน ServiceMonitor แล้วเขียน `PrometheusRule` 2 ตัว: เตือนเมื่อ cert เหลือน้อยกว่า
  14 วัน (cert-manager ต่ออายุตอนเหลือ 30 วัน alert จึงดังเฉพาะเมื่อการต่ออายุล้มเหลว
  ต่อเนื่อง) และเมื่อ cert ไม่ Ready นาน 15 นาที Alertmanager ส่งเข้า Discord เฉพาะ
  alert เรื่อง cert (กัน alert ของ component ที่ k3s ไม่ได้แยก process ออกมา) และแจ้ง
  ตอนปัญหาหายด้วย ทดสอบทั้งเส้นโดยลดเงื่อนไขชั่วคราวให้ alert ดังจริง webhook URL
  เก็บใน Ansible Vault

---

## CI/CD และ Branching Strategy

โปรเจคใช้ flow **dev → main** พร้อม branch protection:

1. งานทั้งหมดทำบน branch **`dev`**
2. ทุกครั้งที่ push/PR **CI workflow** จะรันอัตโนมัติ 2 job:
   - `test` — ตรวจ `go mod tidy` → `go vet` → `go build` → `go test`
   - `security` — Trivy สแกน repository และ container image
3. **Build image ครั้งเดียว** เป็น `devopsminiproject-dev:<tag>` แล้ว deploy ขึ้น
   **dev** เพื่อทดสอบบน cluster จริง
4. **Pull Request** จาก `dev` เข้า `main` ต้องให้ `test` และ `security` ผ่านก่อน
   ถึง merge ได้ (บังคับด้วย branch ruleset บน `main`)
5. **Promote** image tag เดิมจาก `-dev` ไป `-prd` (รันจาก `main` เท่านั้น) โดยไม่
   build ใหม่ และตรวจว่า digest ตรงกัน
6. **Deploy ขึ้น production แบบ manual** ด้วย tag เดิม (คุมจังหวะการ release เอง)

| Workflow | Trigger | หน้าที่ |
|---|---|---|
| `ci.yaml` | push / PR | test (tidy, vet, build, test) และ security (Trivy scan) |
| `workflow.yaml` | manual | build image **dev** แล้ว push เข้า GHCR (Docker Buildx) ใช้ tag แบบ calendar versioning |
| `promote.yaml` | manual | promote image จาก dev ไป prd โดยไม่ build ใหม่ และตรวจ digest |
| `workflow-deploy.yaml` | manual | deploy ขึ้น dev หรือ prd (self-hosted runner แยก label ต่อ environment) |

> ใช้ 2 package ใน GHCR (`devopsminiproject-dev` และ `devopsminiproject-prd`) แทน
> registry แยกต่อ environment ใน `-prd` จึงมีเฉพาะ image ที่ผ่านการทดสอบบน dev และ
> ผ่าน PR มาแล้วเท่านั้น

> ตอนนี้เป็น **Continuous Delivery** — pipeline พร้อม deploy ตลอด แต่ขั้นขึ้น
> production ตั้งใจให้เป็น manual เพื่อคุมจังหวะการ release ส่วน CI gate
> ทำหน้าที่ปกป้อง branch `main`

---

## โครงสร้างโปรเจค (Repository Structure)

```
.
├── cmd/                      # Entrypoint ของแอป (โหมด server / migrator)
├── di/                       # Dependency injection: config, database, server
├── service/                  # HTTP handlers (user CRUD, status, health) + unit tests
├── repository/               # Data access layer (GORM)
├── entity/                   # Domain models
├── util/                     # Build info, helpers + unit tests
├── k8s/
│   ├── platform/             # ตั้งครั้งเดียวต่อ cluster (ไม่ได้ deploy ทุก release)
│   │   ├── cluster-issuer.yaml       # Let's Encrypt staging + production
│   │   ├── cert-manager-values.yaml  # ค่าของ cert-manager (เปิด ServiceMonitor)
│   │   ├── cert-alerts.yaml          # PrometheusRule: cert ใกล้หมดอายุ / ไม่ Ready
│   │   └── monitoring-values.yaml    # ค่าของ kube-prometheus-stack + Alertmanager route
│   ├── namespace.yaml
│   ├── configmap-dev.yaml / configmap-prd.yaml
│   ├── app.yaml              # Deployment + Service
│   ├── job.yaml              # DB migration Job
│   ├── hpa.yaml              # HorizontalPodAutoscaler
│   ├── ingress-dev.yaml      # HTTP
│   └── ingress-prd.yaml      # HTTPS + cert-manager + redirect middleware
├── ansible/                  # Provision k3s cluster (IaC)
│   ├── inventory.ini         # VM ทดสอบ (master / workers + k3s_version)
│   ├── inventory-prd.ini     # VPS production
│   ├── site.yaml             # playbook หลัก: common → k3s_server → k3s_agent → platform
│   ├── requirements.yaml     # pin Ansible collections
│   ├── secrets/              # vm.yaml / prd.yaml เข้ารหัสด้วย Ansible Vault
│   └── roles/
│       ├── common/           # SSH keys, SSH hardening, hostname, cloud-init, netplan, ufw
│       ├── k3s_server/       # ติดตั้ง server, อ่าน token, ตั้ง kubeconfig
│       ├── k3s_agent/        # ติดตั้ง agent แล้ว join master
│       └── platform/         # helm, kube-prometheus-stack, cert-manager, ClusterIssuer, alert rules
├── .github/workflows/        # CI + security scan, build, promote, deploy pipelines
└── Dockerfile                # Multi-stage build
```

---

## API Endpoints

| Method | Path | คำอธิบาย |
|---|---|---|
| `GET` | `/` | ภาพรวมสถานะ runtime |
| `GET` | `/about` | ข้อมูล service (version, build time, commit) |
| `GET` | `/livez` | Liveness probe (เช็คแค่ process) |
| `GET` | `/readyz` | Readiness probe (เช็ค DB) |
| `GET` | `/user` | ดูรายชื่อ user ทั้งหมด |
| `POST` | `/user` | สร้าง user |
| `DELETE` | `/user` | ลบ user |

---

## Monitoring และ Autoscaling

- **Grafana dashboards** (จาก kube-prometheus-stack) ติดตาม CPU/memory ราย pod,
  ราย namespace และราย node
- **Alerting** — แจ้งเตือนเข้า Discord เมื่อ TLS certificate ใกล้หมดอายุหรือไม่ Ready
- **Load testing ด้วย k6** ยิง traffic เข้า API เพื่อยืนยันว่า HPA scale replicas
  ขึ้นเมื่อ CPU สูงต่อเนื่อง และ scale ลงหลังพ้น stabilization window

<!-- ย้ายไฟล์ image.png เดิมมาไว้ที่ docs/ แล้วตั้งชื่อให้ตรงกับด้านล่าง -->
![Grafana dashboard แสดง CPU และ memory ของ pod](docs/grafana.png)
![HPA scale replicas ระหว่าง load test ด้วย k6](docs/hpa-scaling.png)

---

## การตั้ง cluster ใหม่ (Platform Setup)

ขั้นตอนที่ทำครั้งเดียวต่อ cluster ก่อน deploy แอปผ่าน workflow ครั้งแรก:

1. เตรียมเครื่อง Ubuntu 24.04 ที่มี **IP คงที่** และ SSH ด้วย key ได้ แล้วใส่ใน
   `ansible/inventory.ini`
2. สร้างไฟล์ secret ด้วย Ansible Vault (ใส่ `discord_webhook_url`) และติดตั้ง collection:
   ```bash
   cd ansible
   ansible-galaxy collection install -r requirements.yaml
   ansible-vault create secrets/vm.yaml
   ```
3. สร้าง cluster และ platform ทั้งหมดในคำสั่งเดียว (OS, firewall, k3s server + agent,
   helm, kube-prometheus-stack, cert-manager, ClusterIssuer, alert rules):
   ```bash
   ansible-playbook site.yaml -K --ask-vault-pass
   # prd: ansible-playbook -i inventory-prd.ini site.yaml --ask-vault-pass --check --diff
   ```
4. ติดตั้ง self-hosted runner บน node ที่มี kubeconfig และตั้ง label ตาม environment
5. ถ้าผู้ให้บริการมี firewall ภายนอก ให้เปิด port 22, 80 และ 443 (firewall ของ OS
   จัดการโดย Ansible แล้ว — เปิดสาธารณะแค่ 3 port นี้ และให้ node คุยกันได้ทุก port)

> cluster prd ถูกตั้งด้วยมือก่อนมี playbook ขั้นตอนที่ทำมือและปัญหาที่เจอ
> (hostname ซ้ำ, search domain, kubeconfig) ถูกนำมาเขียนเป็น playbook ทดสอบบน VM
> แล้วนำมาใช้กับ prd ด้วย `ansible-playbook -i inventory-prd.ini site.yaml --check --diff`
> ก่อน apply จริง ส่วน Helm charts ที่เคยติดตั้งด้วยมือถูกเก็บ values ไว้ใน git ตั้งแต่แรก
> playbook จึงรับมาดูแลต่อได้โดยไม่มีการเปลี่ยนแปลง ปัจจุบัน prd ถูกจัดการด้วย playbook
> ชุดเดียวกันทั้งหมด (รันซ้ำได้ `changed=0`)

---

## ผู้จัดทำ (Author)

- **Name:** Thanut Sukprasertsom
- **Email:** thanutsukprasertsomm@hotmail.com
- **GitHub:** https://github.com/thanutgit
- **LinkedIn:** https://www.linkedin.com/in/thanut-sukprasertsom-b77048386/