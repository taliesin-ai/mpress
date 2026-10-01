# Resource-limit quality and contract review

The shared budget is local to one load and spans every mounted bundle. Per-bundle
accounting resets for each snapshot. Stored reads and decompression remain
bounded after a size check, so a file growing after stat cannot evade the cap.
A plain file size is only an allocation hint; shrink/growth fixtures cover the
EOF and excess-byte paths. Gzip checksums and trailers are read on accepted
streams, including multiple members. Input bytes are bounded, not peak RSS.

The manifest validation and replacement writer were extracted after initial
complexity findings. Final source metrics introduce no complexity/length
regression. Eight dead-code rows are Go test entry points executed by recorded
runs. Two churn gate rows, LoadAllRoot and loadVersionSite, reflect consecutive
confinement, identity and resource fixes in the same loaders. The signature and
call review confirms why they changed: one shared budget must flow through
current and version loads. Returning to separate per-version budgets to remove
churn would lose the aggregate safety contract. The remaining loadSite churn
row is minor. These findings are retained and reviewed without acknowledgements;
the working-tree quality command therefore still exits 2.

Earlier source-overlay replays exclude the later budget test file (and, for
pre-rooted source, its implementation) because those private seams did not
exist at their baseline. The new budget replay explicitly forwards those seams
to original unbounded production readers; it reproduces small-fixture limit
failures without compilation failures or memory exhaustion. Controls preserve
old truncation/checksum, schema, digest, exact-boundary and replacement behavior.
