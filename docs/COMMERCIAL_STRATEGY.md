# Commercialization & Integration Strategy

## 1. Sales Strategy: Standalone vs. Suite

You have built three modular tools. **Yes, you can and should sell them independently**, but the "Enterprise Suite" is where the big contracts are.

### Product A: Dependency-CI (The "Trojan Horse")
*   **Value Prop:** "Cut your AWS/GitHub Actions bill by 50% instantly."
*   **Target:** Engineering Managers / CFOs.
*   **Sales Model:** **Standalone**.
    *   Free Tier: Open source CLI for local use.
    *   Paid: $20/mo/user for the "Enterprise Report" (Audit logs + Flaky test analytics).
*   **Why Independent?** It's the easiest to adopt. No infrastructure changes needed. It gets your foot in the door.

### Product B: LivePatch (The "Developer Delight")
*   **Value Prop:** "Stop waiting 5 minutes for Docker builds. See changes in 2 seconds."
*   **Target:** Individual Developers / Team Leads.
*   **Sales Model:** **Standalone**.
    *   Free: Localhost syncing.
    *   Paid: $50/mo/dev for **Remote Kubernetes Syncing** (Secure Tunneling).
*   **Why Independent?** Solves a specific pain point for K8s users. Non-K8s users won't care.

### Product C: Quantum Merge (The "Process Fixer")
*   **Value Prop:** "Never break main again. Merge 10x faster."
*   **Target:** VP of Engineering / CTOs of large teams (50+ devs).
*   **Sales Model:** **Bundle Only (Recommended)** or High-Ticket Standalone.
    *   This requires setup (webhooks, server). Usually sold as a managed SaaS ($500+/mo).
*   **Why Bundle?** It works best when combined with Dependency-CI to make the speculative builds fast.

---

## 2. Integration Guide: "Plug & Play"

Here is how you integrate these tools into existing systems without rewriting code.

### A. Integrating Dependency-CI
**Scenario:** A company uses **GitHub Actions** to run tests.

**Current Workflow:**
```yaml
- run: npm install
- run: npm test  # <-- Slow! Runs everything.
```

**New Workflow (Integration):**
```yaml
- run: npm install
- run: curl -L https://your-site.com/dependency-ci -o dep-ci && chmod +x dep-ci
# One line change!
- run: ./dep-ci run --cmd "npm test" --files "${{ github.event.pull_request.changed_files }}"
```
*   **Friction:** Near Zero.
*   **Selling Point:** "Add one line to your YAML, save 20 minutes per build."

### B. Integrating LivePatch
**Scenario:** A team uses **Kubernetes** for their dev environment.

**Step 1: The Dockerfile Change (One time)**
```dockerfile
# Add the agent to the image
COPY --from=velocity-trinity/agent:latest /agent /usr/local/bin/agent
# Run it in background
CMD ["/usr/local/bin/agent", "&", "node", "server.js"]
```

**Step 2: The Developer Workflow**
The developer doesn't use `kubectl`. They just run:
```bash
live-patch sync ./src --target my-dev-pod.namespace:8080 --restart "npm restart"
```
*   **Friction:** Low. Requires adding the binary to the base image.

### C. Integrating Quantum Merge
**Scenario:** Replacing the standard GitHub Merge Button.

1.  **Disable Direct Merge:** go to GitHub Settings -> Branch Protection -> "Require status checks to pass".
2.  **Add Webhook:** Point GitHub Webhooks to `https://your-quantum-instance.com/webhook`.
3.  **The Trigger:**
    *   Developer comments `/merge` on the PR.
    *   **Quantum Merge** catches the webhook.
    *   It triggers a **Jenkins Job** (via Jenkins API) or **GitHub Action** (via `workflow_dispatch`).
    *   If it passes, Quantum Merge calls GitHub API to `PUT /pulls/:id/merge`.

---

## 3. The "Velocity Trinity" Bundle Pricing

| Tier | Price | Includes |
| :--- | :--- | :--- |
| **Starter** | Free | • Dependency-CI (Local)<br>• LivePatch (Localhost) |
| **Team** | $49/dev/mo | • **Dependency-CI (CI/CD)**<br>• **LivePatch (Remote K8s)**<br>• Email Support |
| **Enterprise** | Custom | • **Quantum Merge (Managed)**<br>• On-Prem Deployment<br>• SSO / Audit Logs |

**Recommendation:** Start by selling **Dependency-CI** to get trusted by the team. Then upsell **LivePatch** to their devs. Finally, sell **Quantum Merge** to their CTO to fix their process.
