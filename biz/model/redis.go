package model

const (
	ARTICLE_POOL        = "ArticlePool_"       //不同类型文章池
	FRIEND_ARTICLE_LIST = "FriendArticleList_" //朋友文章列表
	FOLLOW_ARTICLE_LIST = "FollowArticleList_" //关注文章列表
	GIVELIKE            = "give_like_"         //点赞文章列表
	COMMENT_GIVELIKE    = "comment:giveLike:"  //评论点赞列表
	NOTIFY_HISTORY      = "NotifyHistory:"     //历史通知
	UNREAD_NOTIFY       = "UnreadNotify:"      //未读通知数目
	BLOOM_FILTER        = "bloomFilter_"       //布隆过滤器
	RELATION            = "relation_"          //用户之间关系
	MESSAGE_HISTORY     = "MessageHistory_"    //聊天历史消息
	UNREAD_CHAT         = "UnreadChat:"        //聊天未读消息数目
	SEARCH_RECORD       = "SearchRecord:"      //搜索记录
)
