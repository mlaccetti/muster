# Cursor Agent projection

Muster cannot configure Cursor's global User Rules through a supported CLI or file interface. Add the contents of Muster's `policy/orchestrator.md` to Cursor Settings → Rules → User Rules to make Muster the default orchestrator for Cursor.

Cursor Agent resolves as `cursor-agent`. Worker delegation and event-driven waiting depend on capabilities exposed by the active Cursor session; the projection policy asks the orchestrator to use only capabilities actually available.
