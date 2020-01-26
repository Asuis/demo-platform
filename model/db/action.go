package db

type ActionType int

const (
	ActionCreateRepo           ActionType = iota + 1 // 1
	ActionRenameRepo                                 // 2
	ActionStarRepo                                   // 3
	ActionWatchRepo                                  // 4
	ActionCommitRepo                                 // 5
	ActionCreateIssue                                // 6
	ActionCreatePullRequest                          // 7
	ActionTransferRepo                               // 8
	ActionPushTag                                    // 9
	ActionCommentIssue                               // 10
	ActionMergePullRequest                           // 11
	ActionCloseIssue                                 // 12
	ActionReopenIssue                                // 13
	ActionClosePullRequest                           // 14
	ActionReopenPullRequest                          // 15
	ActionCreateBranch                               // 16
	ActionDeleteBranch                               // 17
	ActionDeleteTag                                  // 18
	ActionForkRepo                                   // 19
	ActionMirrorSyncPush                             // 20
	ActionMirrorSyncCreate                           // 21
	ActionMirrorSyncDelete                           // 22
)
