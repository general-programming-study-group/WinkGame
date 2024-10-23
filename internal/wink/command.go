package wink

import (
  // system packages
  "log"

  // internal packages
  "github.com/coex1/EchoBot/internal/data"

  // external packages
  dgo "github.com/bwmarrin/discordgo"
)


func CommandHandle(s *dgo.Session, event *dgo.InteractionCreate, guild *data.Guild) {
	var minListCnt int = MIN_PLAYER_CNT
	var err error
	var optionList []dgo.SelectMenuOption
	var members []*dgo.Member

	guild.Wink.SelectedUsers[event.GuildID] = make([]string, 0)

	guild.Wink.CheckedUsers = make(map[string]bool)
	guild.Wink.TotalParticipants = 0
	guild.Wink.MessageIDMap = make(map[string]string)

	// get guild members
	members, err = s.GuildMembers(event.GuildID, QUERY_STRING, MAX_MEMBER_GET)
	if err != nil {
		log.Fatalf("Failed getting members [%v]", err)
		return
	}

	// create select list from 'members'
	for _, m := range members {
		// check if 'm' is a bot
		if m.User.Bot {
			continue
		}

		optionList = append(optionList, dgo.SelectMenuOption{
			Label: m.User.GlobalName,
			Value: m.User.ID,
		})
	}

  response := &dgo.InteractionResponse{
		Type: dgo.InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
      Embeds: []*dgo.MessageEmbed{ 
        {
          Title:        "게임 참여자 선택!",
          Description:  "게임에 참석할 사용자들을 선택해 주세요."+
                        "\n최소 5명 이상이 선택 되어야 게임이 가능합니다."+
                        "\n선택 하셨으면 '게임시작' 버튼을 클릭 해 주세요.",
          Color:        0x2AFF00,
        },
      },
			Components: []dgo.MessageComponent{
        dgo.ActionsRow{
          Components: []dgo.MessageComponent{
            dgo.SelectMenu{
              CustomID:     "wink_Start_listUpdate",
              Placeholder:  "사용자 목록",
              MinValues:    &minListCnt,
              MaxValues:    len(optionList),
              Options:      optionList,
            },
          },
        },
        dgo.ActionsRow{
          Components: []dgo.MessageComponent{
            &dgo.Button{
              Label:    "게임시작",          // 버튼 텍스트
              Style:    dgo.PrimaryButton,   // 버튼 스타일
              CustomID: "wink_Start_Button", // 버튼 클릭 시 처리할 ID
            },
          },
        },
			},
		},
	}

  // respond to command by sending Start Menu
	err = s.InteractionRespond(event.Interaction, response)
	if err != nil {
		log.Fatalf("Failed to send response [%v]", err)
		return
	}
}
