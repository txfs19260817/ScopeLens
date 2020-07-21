<template>
  <v-container class="about">
    <v-row>
      <v-col>
        <h1>欢迎访问ScopeLens宝可梦队伍分享平台</h1>
        <h2>基本功能</h2>
        <p>您可以在本网站上与全世界的玩家们分享您认为不错的宝可梦队伍。</p>
        <h3>分享队伍</h3>
        <p>
          您不仅可以选择分享游戏软件中导出的租借队伍ID界面，
          还可以上传<a href="https://play.pokemonshowdown.com/">Pokemon Showdown</a>格式的队伍文本（下面简称Showdown文本）。
          <strong>
            我们强烈建议您在发布队伍时同时上传两者。请确保上传的图片包含租借ID。只上传Showdown文本的话宝可梦名字不要带昵称。
          </strong>
          如果您选择仅分享Showdown文本，本站会从文本导出一张仿照游戏中的租借队预览图，
          记录队伍基本信息和标题与作者，代替租借ID显示在我们的网站上。
        </p>
        <h3>检索队伍</h3>
        <p>可通过指定对战模式和一只或复数只宝可梦对服务器存储的队伍进行检索，可为组队过程提供参考。</p>
        <h3>使用率查看</h3>
        <p>
          可通过指定对战模式查看服务器存储的所有队伍下的宝可梦使用率。
          <strong>注意：若指定的模式下没有任一种类宝可梦使用频次大于1，则不会返回任何结果。</strong>
        </p>
        <h2>使用须知</h2>
        <p>
          <ol>
            <li>提交队伍需<strong>登录</strong></li>
            <li>
              本站对POST请求（包括上传队伍和队伍检索）有频度限制，如果您遇见
              <i>"Sorry, you are visiting our service too frequent. Please try again after a minute. (429)"</i>
              这样的错误提示，请稍等1分钟后重试。
            </li>
            <li>
              忘记密码<del>建议重新注册反正也不是很珍贵</del>
              请用您注册时的邮箱发送邮件（见下方联系方式）到我们的邮箱来重置密码。
            </li>
          </ol>
        </p>
        <h2>联系方式</h2>
        <p>
          如果您在使用本网站中遇到任何问题，或愿意作出贡献，
          热烈欢迎到<a href="https://github.com/txfs19260817/ScopeLens">Github</a>上提交issue或pull request，
          或者发送<a href="mailto:zeminjiang@mail.ru?subject=ScopeLens">邮件</a>进行反馈。
        </p>
        <p>
          如果您愿意为本站添加其他语种翻译，请参阅
          <a href="https://github.com/txfs19260817/ScopeLens/tree/master/web/scopelens-website/src/locales">该目录</a>
          下的文件进行翻译，并PR到该项目上，提前感谢您的贡献。
        </p>
        <p>唯一指定只聊双打群QQ群：526116328</p>
        <h2>Acknowledgement</h2>
        <ul>
          <li>本网站使用的资源（图标等）版权归 © Nintendo/Creatures Inc./GAME FREAK Inc. 所有</li>
          <li>本站UI借鉴自项目<a href="https://github.com/OpenEpicData/FlamingoWeb">FlamingoWeb</a></li>
        </ul>
      </v-col>
    </v-row>
  </v-container>
</template>
