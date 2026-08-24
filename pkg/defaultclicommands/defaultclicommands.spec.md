# defaultclicommands specifications

##  Implementation

- This package and all its subpackages are ment to wire up the commands with the actual implementation.
    - Do not add additional functions or logic in the `defaultclicommands` package or its subpackages.
    - If functionality is missing add the implementation a the corresponding subpackage of `pkg`.
- Add a usage example in the `Long` description for every `cobra.Command`.
    - Begin the `Long` description by reusing the `short` one:
        - Example:
            ```golang
            func NewConcatFilesToKnowledgeFileCmd() *cobra.Command {
                const short = "Concat files in a directory to one knowledge file."

                cmd := &cobra.Command{
                    Use:   "concat-files-to-knowledge-file",
                    Short: short,
                    Long: short + `

            [Additional description, links, hints]

            Usage:
                ` + os.Args[0] + ` ai concat-files-to-knowledge-file --verbose [toplevel dir with knowledgefiles] > documentation.markdown

            [Additional Usage examples]
            `,
                    Run: func(cmd *cobra.Command, args []string) {
                        // ...
                    },
                }

                return cmd
            }
            ```
