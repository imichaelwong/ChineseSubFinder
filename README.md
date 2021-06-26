# ChineseSubFinder

本项目的初衷仅仅是想自动化搞定**限定条件**下 **中文** 字幕下载。

> 开发中，可能有不兼容性的调整（配置文件字段变更）
>
> 最新版本 v0.8.x 配置文件 config.xml 有更新，注意看下面的文档
>
> v0.7.x 优化视频跳过下载字幕的逻辑，加快扫描速度

## Why？

注意，因为近期参考《[高阶教程-追剧全流程自动化 | sleele的博客](https://sleele.com/tag/高阶教程-追剧全流程自动化/)》搞定了自动下载，美剧、电影没啥问题。但是遇到字幕下载的困难，里面推荐的都不好用，能下载一部分，大部分都不行。当然有可能是个人的问题。为此就打算自己整一个专用的下载器。

手动去下载再丢过去改名也不是不行，这不是懒嘛...

首先，明确一点，因为搞定了 sonarr 和 radarr 以及 Emby，同时部分手动下载的视频也会使用 tinyMediaManager 去处理，所以可以认为所有的视频是都有 IMDB ID 的。那么就可以取巧，用 IMDB ID 去搜索（最差也能用标准的视频文件名称去搜索嘛）。

## 功能

### 支持的部署方式

* docker，见 [How to use](https://github.com/allanpk716/ChineseSubFinder/blob/master/DesignFile/HowToUse.md)

### 支持的视频分类（如何削刮视频的）

|      类型       | Emby/Jellyfin | TinyMediaManager | Sonarr | Radarr | 人工随意命名分类 |                             备注                             |
| :-------------: | :-----------: | :--------------: | :----: | :----: | :--------------: | :----------------------------------------------------------: |
|      电影       |       ✔       |        ✔         |   ✖    |   ✔    |        ✓         |              通过 IMDB ID 或者 文件名 进行搜索               |
|     连续剧      |       ✔       |        ✔         |   ✔    |   ✖    |        ✖         |             **必须**依赖 tvshow.nfo 中的 IMDB ID             |
| 日本动画(Anime) |       ✖       |        ✖         |   ✖    |   ✖    |        ✖         | [待定，见讨论](https://github.com/allanpk716/ChineseSubFinder/issues/1) |

* ✔ -- 支持
* ✓ -- 支持，但是可能搜索结果不稳定
* ✖ -- 不支持

### 支持的字幕下载站点

| 字幕站点 | 电影 | 连续剧 | Anime |
| :------: | :--: | :----: | :---: |
|  zimuku  |  ✔   |   ✔    |   ✖   |
|  subhd   |  ✔   |   ✔    |   ✖   |
| shooter  |  ✔   |   ✔    |   ✓   |
|  xunlei  |  ✔   |   ✔    |   ✓   |

| 字幕站点 | 电影目录下有 movie.xml或电影名称.nfo | 连续剧目录下有 tvshow.nfo | 通过视频唯一ID | 视频文件名 |
| :------: | :----------------------------------: | :-----------------------: | :------------: | :--------: |
|  zimuku  |                  ✔                   |             ✔             |       ✖        |     ✓      |
|  subhd   |                  ✔                   |             ✔             |       ✖        |     ✓      |
| shooter  |                  ✖                   |             ✖             |       ✔        |     ✖      |
|  xunlei  |                  ✖                   |             ✖             |       ✔        |     ✖      |

* ✔ -- 支持
* ✓ -- 支持，但是可能搜索结果不稳定
* ✖ -- 不支持

### 支持的视频格式

* mp4
* mkv
* rmvb
* iso

### 字幕网站优先级

网站字幕优先级别暂定 ：zimuku -> subhd -> xunlei -> shooter，暂时不支持修改优先级

### 字幕格式优先级

| 字幕格式的优先级选择 | 根据网站和字幕的排名自动选择（字幕类型不定） | 优先 srt | 优先 ass/ssa |
| :------------------: | :------------------------------------------: | :------: | :----------: |
|   SubTypePriority    |                      0                       |    1     |      2       |

### 字幕语言类型优先级

* 双语 -> 单语种
* 简体 -> 繁体

### 支持字幕语言的检测

并非简单的从文件名或者字幕提供方的标记进行识别，而是读取字幕文件，进行双语、简体、繁体等的识别，支持列表如下：

```go
const (
	Unknow                     Language = iota // 未知语言
	ChineseSimple                              // 简体中文
	ChineseTraditional                         // 繁体中文
	ChineseSimpleEnglish                       // 简英双语字幕
	ChineseTraditionalEnglish                  // 繁英双语字幕
	English                                    // 英文
	Japanese                                   // 日语
	ChineseSimpleJapanese                      // 简日双语字幕
	ChineseTraditionalJapanese                 // 繁日双语字幕
	Korean                                     // 韩语
	ChineseSimpleKorean                        // 简韩双语字幕
	ChineseTraditionalKorean                   // 繁韩双语字幕
)
```

然后相应的会转换为以下的字幕语言“后缀名”

```go
// 需要符合 emby 的格式要求，在后缀名前面
const (
	Emby_unknow = ".unknow"					// 未知语言
	Emby_chs 	= ".chs"					// 简体
	Emby_cht 	= ".cht"					// 繁体
	Emby_chs_en = ".chs_en"                 // 简英双语字幕
	Emby_cht_en = ".cht_en"                	// 繁英双语字幕
	Emby_en 	= ".en"                       // 英文
	Emby_jp 	= ".jp"						// 日语
	Emby_chs_jp = ".chs_jp"                 // 简日双语字幕
	Emby_cht_jp = ".cht_jp"                	// 繁日双语字幕
	Emby_kr 	= ".kr"                     // 韩语
	Emby_chs_kr = ".chs_kr"                 // 简韩双语字幕
	Emby_cht_kr = ".cht_kr"                	// 繁韩双语字幕
)
```


## How to use

[How To Use](https://github.com/allanpk716/ChineseSubFinder/blob/master/DesignFile/HowToUse.md)

## 其他文档

* [削刮器的推荐设置](https://github.com/allanpk716/ChineseSubFinder/blob/master/DesignFile/%E5%89%8A%E5%88%AE%E5%99%A8%E7%9A%84%E6%8E%A8%E8%8D%90%E8%AE%BE%E7%BD%AE.md)
* [如何手动刷新 Emby 加载字幕](https://github.com/allanpk716/ChineseSubFinder/blob/master/DesignFile/%E5%A6%82%E4%BD%95%E6%89%8B%E5%8A%A8%E5%88%B7%E6%96%B0%20Emby%20%E5%8A%A0%E8%BD%BD%E5%AD%97%E5%B9%95.md)
* [连续剧如何搜索字幕](https://github.com/allanpk716/ChineseSubFinder/blob/master/DesignFile/%E8%BF%9E%E7%BB%AD%E5%89%A7%E5%A6%82%E4%BD%95%E6%90%9C%E7%B4%A2%E5%AD%97%E5%B9%95.md)
* [设计](https://github.com/allanpk716/ChineseSubFinder/blob/master/DesignFile/%E8%AE%BE%E8%AE%A1.md)

## 版本

* v0.8.x 调整 docker 镜像结构 -- 2021年6月25日
* v0.7.x 提高搜索效率 -- 2021年6月25日
* v0.6.x 支持设置字幕格式的优先级 -- 2021年6月23日
* v0.5.x 支持连续剧字幕下载 -- 2021年6月19日
* v0.4.x 支持设置并发数 -- 2021年6月18日
* v0.3.x 支持连续剧字幕下载（连续剧暂时不支持 subhd） -- 2021年6月17日
* v0.2.0 docker 版本支持 subhd 的下载了，镜像体积也变大了 -- 2021年6月14日
* 完成初版，仅仅支持电影的字幕下载 -- 2021年6月13日

## TODO

* 字幕的风评（有些字幕太差了，需要进行过滤，考虑排除，字幕组，关键词，机翻，以及评分等条件
* 加入 Web 设置界面（也许没得很大的必要···）
* 提供 API 接口，部署后，允许额外的程序访问（类似 emby 等）获取字幕
* 支持 Anime 的字幕下载

## 感谢

感谢下面项目的帮助

* [Andyfoo/GoSubTitleSearcher: 字幕搜索查询(go语言版)](https://github.com/Andyfoo/GoSubTitleSearcher)
* [go-rod/rod: A Devtools driver for web automation and scraping (github.com)](https://github.com/go-rod/rod)
* [ausaki/subfinder: 字幕查找器 (github.com)](https://github.com/ausaki/subfinder)
* [golandscape/sat: 高性能简繁体转换 (github.com)](https://github.com/golandscape/sat)


# 预览图
![Xnip2021-06-25_11-11-55](https://cdn.jsdelivr.net/gh/SuperNG6/pic@master/uPic/2021-06-25/Xnip2021-06-25_11-11-55.jpg)
![Xnip2021-06-25_11-12-33](https://cdn.jsdelivr.net/gh/SuperNG6/pic@master/uPic/2021-06-25/Xnip2021-06-25_11-12-33.jpg)
![Xnip2021-06-25_10-29-06](https://cdn.jsdelivr.net/gh/SuperNG6/pic@master/uPic/2021-06-25/Xnip2021-06-25_10-29-06.jpg)
![Xnip2021-06-25_10-24-22](https://cdn.jsdelivr.net/gh/SuperNG6/pic@master/uPic/2021-06-25/Xnip2021-06-25_10-24-22.jpg)
![Xnip2021-06-25_11-42-38](https://cdn.jsdelivr.net/gh/SuperNG6/pic@master/uPic/2021-06-25/Xnip2021-06-25_11-42-38.jpg)