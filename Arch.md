🧠 High-Level Architecture Reminder
Your project "KubeVision" has 4 major components:
Backend API (FastAPI)
— receives metrics from agent
— runs ML
— exposes recommendations
— stores data
In-cluster Agent (Go)
— runs inside Kubernetes
— collects metrics
— queries K8s API
— sends data to backend
ML Engine
— forecasting, rightsizing, efficiency score
CLI tool
— for developers to query insights
Therefore your repo must support all 4 cleanly.