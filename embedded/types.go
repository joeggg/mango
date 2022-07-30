package embedded

import (
	"fmt"
	"mango/pb"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

type EmbeddedHandler func(proto.Message) error

var EmbdeddedHandlers = map[int][]EmbeddedHandler{}

// Embedded type to proto struct name
var EmbeddedTypes = map[int]string{
	int(pb.NET_Messages_net_NOP):                        "mango.CNETMsg_NOP",
	int(pb.NET_Messages_net_Disconnect):                 "mango.CNETMsg_Disconnect",
	int(pb.NET_Messages_net_Tick):                       "mango.CNETMsg_Tick",
	int(pb.NET_Messages_net_SetConVar):                  "mango.CNETMsg_SetConVar",
	int(pb.NET_Messages_net_SignonState):                "mango.CNETMsg_SignonState",
	int(pb.NET_Messages_net_SpawnGroup_Load):            "mango.CNETMsg_SpawnGroup_Load",
	int(pb.NET_Messages_net_SpawnGroup_ManifestUpdate):  "mango.CNETMsg_SpawnGroup_ManifestUpdate",
	int(pb.NET_Messages_net_SpawnGroup_SetCreationTick): "mango.CNETMsg_SpawnGroup_SetCreationTick",
	int(pb.NET_Messages_net_SpawnGroup_Unload):          "mango.CNETMsg_SpawnGroup_Unload",
	int(pb.NET_Messages_net_SpawnGroup_LoadCompleted):   "mango.CNETMsg_SpawnGroup_LoadCompleted",

	int(pb.CLC_Messages_clc_ClientInfo):      "mango.CCLCMsg_ClientInfo",
	int(pb.CLC_Messages_clc_Move):            "mango.CCLCMsg_Move",
	int(pb.CLC_Messages_clc_BaselineAck):     "mango.CCLCMsg_BaselineAck",
	int(pb.CLC_Messages_clc_LoadingProgress): "mango.CCLCMsg_LoadingProgress",

	int(pb.SVC_Messages_svc_ServerInfo):           "mango.CSVCMsg_ServerInfo",
	int(pb.SVC_Messages_svc_FlattenedSerializer):  "mango.CSVCMsg_FlattenedSerializer",
	int(pb.SVC_Messages_svc_ClassInfo):            "mango.CSVCMsg_ClassInfo",
	int(pb.SVC_Messages_svc_CreateStringTable):    "mango.CSVCMsg_CreateStringTable",
	int(pb.SVC_Messages_svc_UpdateStringTable):    "mango.CSVCMsg_UpdateStringTable",
	int(pb.SVC_Messages_svc_VoiceInit):            "mango.CSVCMsg_VoiceInit",
	int(pb.SVC_Messages_svc_VoiceData):            "mango.CSVCMsg_VoiceData",
	int(pb.SVC_Messages_svc_Print):                "mango.CSVCMsg_Print",
	int(pb.SVC_Messages_svc_SetView):              "mango.CSVCMsg_SetView",
	int(pb.SVC_Messages_svc_ClearAllStringTables): "mango.CSVCMsg_ClearAllStringTables",
	int(pb.SVC_Messages_svc_PacketEntities):       "mango.CSVCMsg_PacketEntities",
	int(pb.SVC_Messages_svc_PeerList):             "mango.CSVCMsg_PeerList",
	int(pb.SVC_Messages_svc_HLTVStatus):           "mango.CSVCMsg_HLTVStatus",
	int(pb.SVC_Messages_svc_FullFrameSplit):       "mango.CSVCMsg_FullFrameSplit",

	int(pb.EBaseUserMessages_UM_Fade):            "mango.CUserMessageFade",
	int(pb.EBaseUserMessages_UM_ResetHUD):        "mango.CUserMessageResetHUD",
	int(pb.EBaseUserMessages_UM_SayText2):        "mango.CUserMessageSayText2",
	int(pb.EBaseUserMessages_UM_TextMsg):         "mango.CUserMessageTextMsg",
	int(pb.EBaseUserMessages_UM_VoiceMask):       "mango.CUserMessageVoiceMask",
	int(pb.EBaseUserMessages_UM_SendAudio):       "mango.CUserMessageSendAudio",
	int(pb.EBaseUserMessages_UM_AudioParameter):  "mango.CUserMessageAudioParameter",
	int(pb.EBaseUserMessages_UM_ParticleManager): "mango.CDOTAUserMsg_ParticleManager",

	int(pb.EBaseGameEvents_GE_Source1LegacyGameEventList): "mango.CSVCMsg_GameEventList",
	int(pb.EBaseGameEvents_GE_Source1LegacyListenEvents):  "mango.CMsgSource1LegacyListenEvents",
	int(pb.EBaseGameEvents_GE_Source1LegacyGameEvent):     "mango.CSVCMsg_GameEvent",
	int(pb.EBaseGameEvents_GE_SosStartSoundEvent):         "mango.CMsgSosStartSoundEvent",
	int(pb.EBaseGameEvents_GE_SosStopSoundEvent):          "mango.CMsgSosStopSoundEvent",
	int(pb.EBaseGameEvents_GE_SosSetSoundEventParams):     "mango.CMsgSosSetSoundEventParams",
	int(pb.EBaseGameEvents_GE_SosStopSoundEventHash):      "mango.CMsgSosStopSoundEventHash",

	int(pb.ETEProtobufIds_TE_EffectDispatchId): "mango.CMsgTEEffectDispatch",

	int(pb.EDotaUserMessages_DOTA_UM_ChatEvent):                  "mango.CDOTAUserMsg_ChatEvent",
	int(pb.EDotaUserMessages_DOTA_UM_CombatLogBulkData):          "mango.CDOTAUserMsg_CombatLogBulkData",
	int(pb.EDotaUserMessages_DOTA_UM_CreateLinearProjectile):     "mango.CDOTAUserMsg_CreateLinearProjectile",
	int(pb.EDotaUserMessages_DOTA_UM_DestroyLinearProjectile):    "mango.CDOTAUserMsg_DestroyLinearProjectile",
	int(pb.EDotaUserMessages_DOTA_UM_GlobalLightColor):           "mango.CDOTAUserMsg_GlobalLightColor",
	int(pb.EDotaUserMessages_DOTA_UM_GlobalLightDirection):       "mango.CDOTAUserMsg_GlobalLightDirection",
	int(pb.EDotaUserMessages_DOTA_UM_DodgeTrackingProjectiles):   "mango.CDOTAUserMsg_DodgeTrackingProjectiles",
	int(pb.EDotaUserMessages_DOTA_UM_LocationPing):               "mango.CDOTAUserMsg_LocationPing",
	int(pb.EDotaUserMessages_DOTA_UM_MapLine):                    "mango.CDOTAUserMsg_MapLine",
	int(pb.EDotaUserMessages_DOTA_UM_MiniKillCamInfo):            "mango.CDOTAUserMsg_MiniKillCamInfo",
	int(pb.EDotaUserMessages_DOTA_UM_MinimapEvent):               "mango.CDOTAUserMsg_MinimapEvent",
	int(pb.EDotaUserMessages_DOTA_UM_NevermoreRequiem):           "mango.CDOTAUserMsg_NevermoreRequiem",
	int(pb.EDotaUserMessages_DOTA_UM_OverheadEvent):              "mango.CDOTAUserMsg_OverheadEvent",
	int(pb.EDotaUserMessages_DOTA_UM_SharedCooldown):             "mango.CDOTAUserMsg_SharedCooldown",
	int(pb.EDotaUserMessages_DOTA_UM_SpectatorPlayerClick):       "mango.CDOTAUserMsg_SpectatorPlayerClick",
	int(pb.EDotaUserMessages_DOTA_UM_UnitEvent):                  "mango.CDOTAUserMsg_UnitEvent",
	int(pb.EDotaUserMessages_DOTA_UM_ParticleManager):            "mango.CDOTAUserMsg_ParticleManager",
	int(pb.EDotaUserMessages_DOTA_UM_BotChat):                    "mango.CDOTAUserMsg_BotChat",
	int(pb.EDotaUserMessages_DOTA_UM_HudError):                   "mango.CDOTAUserMsg_HudError",
	int(pb.EDotaUserMessages_DOTA_UM_ItemPurchased):              "mango.CDOTAUserMsg_ItemPurchased",
	int(pb.EDotaUserMessages_DOTA_UM_WorldLine):                  "mango.CDOTAUserMsg_WorldLine",
	int(pb.EDotaUserMessages_DOTA_UM_ChatWheel):                  "mango.CDOTAUserMsg_ChatWheel",
	int(pb.EDotaUserMessages_DOTA_UM_GamerulesStateChanged):      "mango.CDOTAUserMsg_GamerulesStateChanged",
	int(pb.EDotaUserMessages_DOTA_UM_SendStatPopup):              "mango.CDOTAUserMsg_SendStatPopup",
	int(pb.EDotaUserMessages_DOTA_UM_SendRoshanPopup):            "mango.CDOTAUserMsg_SendRoshanPopup",
	int(pb.EDotaUserMessages_DOTA_UM_TE_Projectile):              "mango.CDOTAUserMsg_TE_Projectile",
	int(pb.EDotaUserMessages_DOTA_UM_TE_ProjectileLoc):           "mango.CDOTAUserMsg_TE_ProjectileLoc",
	int(pb.EDotaUserMessages_DOTA_UM_TE_DotaBloodImpact):         "mango.CDOTAUserMsg_TE_DotaBloodImpact",
	int(pb.EDotaUserMessages_DOTA_UM_TE_UnitAnimation):           "mango.CDOTAUserMsg_TE_UnitAnimation",
	int(pb.EDotaUserMessages_DOTA_UM_TE_UnitAnimationEnd):        "mango.CDOTAUserMsg_TE_UnitAnimationEnd",
	int(pb.EDotaUserMessages_DOTA_UM_AbilityPing):                "mango.CDOTAUserMsg_AbilityPing",
	int(pb.EDotaUserMessages_DOTA_UM_WillPurchaseAlert):          "mango.CDOTAUserMsg_WillPurchaseAlert",
	int(pb.EDotaUserMessages_DOTA_UM_AbilitySteal):               "mango.CDOTAUserMsg_AbilitySteal",
	int(pb.EDotaUserMessages_DOTA_UM_CourierKilledAlert):         "mango.CDOTAUserMsg_CourierKilledAlert",
	int(pb.EDotaUserMessages_DOTA_UM_EnemyItemAlert):             "mango.CDOTAUserMsg_EnemyItemAlert",
	int(pb.EDotaUserMessages_DOTA_UM_QuickBuyAlert):              "mango.CDOTAUserMsg_QuickBuyAlert",
	int(pb.EDotaUserMessages_DOTA_UM_PredictionResult):           "mango.CDOTAUserMsg_PredictionResult",
	int(pb.EDotaUserMessages_DOTA_UM_ModifierAlert):              "mango.CDOTAUserMsg_ModifierAlert",
	int(pb.EDotaUserMessages_DOTA_UM_HPManaAlert):                "mango.CDOTAUserMsg_HPManaAlert",
	int(pb.EDotaUserMessages_DOTA_UM_SpectatorPlayerUnitOrders):  "mango.CDOTAUserMsg_SpectatorPlayerUnitOrders",
	int(pb.EDotaUserMessages_DOTA_UM_ProjectionAbility):          "mango.CDOTAUserMsg_ProjectionAbility",
	int(pb.EDotaUserMessages_DOTA_UM_ProjectionEvent):            "mango.CDOTAUserMsg_ProjectionEvent",
	int(pb.EDotaUserMessages_DOTA_UM_CombatLogDataHLTV):          "mango.CMsgDOTACombatLogEntry",
	int(pb.EDotaUserMessages_DOTA_UM_XPAlert):                    "mango.CDOTAUserMsg_XPAlert",
	int(pb.EDotaUserMessages_DOTA_UM_UpdateQuestProgress):        "mango.CDOTAUserMsg_UpdateQuestProgress",
	int(pb.EDotaUserMessages_DOTA_UM_MatchMetadata):              "mango.CDOTAMatchMetadataFile",
	int(pb.EDotaUserMessages_DOTA_UM_MatchDetails):               "mango.CMsgDOTAMatch",
	int(pb.EDotaUserMessages_DOTA_UM_SelectPenaltyGold):          "mango.CDOTAUserMsg_SelectPenaltyGold",
	int(pb.EDotaUserMessages_DOTA_UM_RollDiceResult):             "mango.CDOTAUserMsg_RollDiceResult",
	int(pb.EDotaUserMessages_DOTA_UM_FlipCoinResult):             "mango.CDOTAUserMsg_FlipCoinResult",
	int(pb.EDotaUserMessages_DOTA_UM_TeamCaptainChanged):         "mango.CDOTAUserMessage_TeamCaptainChanged",
	int(pb.EDotaUserMessages_DOTA_UM_SendRoshanSpectatorPhase):   "mango.CDOTAUserMsg_SendRoshanSpectatorPhase",
	int(pb.EDotaUserMessages_DOTA_UM_TE_DestroyProjectile):       "mango.CDOTAUserMsg_TE_DestroyProjectile",
	int(pb.EDotaUserMessages_DOTA_UM_HeroRelicProgress):          "mango.CDOTAUserMsg_HeroRelicProgress",
	int(pb.EDotaUserMessages_DOTA_UM_ItemSold):                   "mango.CDOTAUserMsg_ItemSold",
	int(pb.EDotaUserMessages_DOTA_UM_DamageReport):               "mango.CDOTAUserMsg_DamageReport",
	int(pb.EDotaUserMessages_DOTA_UM_SalutePlayer):               "mango.CDOTAUserMsg_SalutePlayer",
	int(pb.EDotaUserMessages_DOTA_UM_TipAlert):                   "mango.CDOTAUserMsg_TipAlert",
	int(pb.EDotaUserMessages_DOTA_UM_EmptyTeleportAlert):         "mango.CDOTAUserMsg_EmptyTeleportAlert",
	int(pb.EDotaUserMessages_DOTA_UM_MarsArenaOfBloodAttack):     "mango.CDOTAUserMsg_MarsArenaOfBloodAttack",
	int(pb.EDotaUserMessages_DOTA_UM_ESArcanaCombo):              "mango.CDOTAUserMsg_ESArcanaCombo",
	int(pb.EDotaUserMessages_DOTA_UM_ESArcanaComboSummary):       "mango.CDOTAUserMsg_ESArcanaComboSummary",
	int(pb.EDotaUserMessages_DOTA_UM_HighFiveLeftHanging):        "mango.CDOTAUserMsg_HighFiveLeftHanging",
	int(pb.EDotaUserMessages_DOTA_UM_HighFiveCompleted):          "mango.CDOTAUserMsg_HighFiveCompleted",
	int(pb.EDotaUserMessages_DOTA_UM_ShovelUnearth):              "mango.CDOTAUserMsg_ShovelUnearth",
	int(pb.EDotaUserMessages_DOTA_EM_InvokerSpellCast):           "mango.CDOTAEntityMsg_InvokerSpellCast",
	int(pb.EDotaUserMessages_DOTA_UM_RadarAlert):                 "mango.CDOTAUserMsg_RadarAlert",
	int(pb.EDotaUserMessages_DOTA_UM_AllStarEvent):               "mango.CDOTAUserMsg_AllStarEvent",
	int(pb.EDotaUserMessages_DOTA_UM_TalentTreeAlert):            "mango.CDOTAUserMsg_TalentTreeAlert",
	int(pb.EDotaUserMessages_DOTA_UM_QueuedOrderRemoved):         "mango.CDOTAUserMsg_QueuedOrderRemoved",
	int(pb.EDotaUserMessages_DOTA_UM_DebugChallenge):             "mango.CDOTAUserMsg_DebugChallenge",
	int(pb.EDotaUserMessages_DOTA_UM_OMArcanaCombo):              "mango.CDOTAUserMsg_OMArcanaCombo",
	int(pb.EDotaUserMessages_DOTA_UM_FoundNeutralItem):           "mango.CDOTAUserMsg_FoundNeutralItem",
	int(pb.EDotaUserMessages_DOTA_UM_OutpostCaptured):            "mango.CDOTAUserMsg_OutpostCaptured",
	int(pb.EDotaUserMessages_DOTA_UM_OutpostGrantedXP):           "mango.CDOTAUserMsg_OutpostGrantedXP",
	int(pb.EDotaUserMessages_DOTA_UM_MoveCameraToUnit):           "mango.CDOTAUserMsg_MoveCameraToUnit",
	int(pb.EDotaUserMessages_DOTA_UM_PauseMinigameData):          "mango.CDOTAUserMsg_PauseMinigameData",
	int(pb.EDotaUserMessages_DOTA_UM_VersusScene_PlayerBehavior): "mango.CDOTAUserMsg_VersusScene_PlayerBehavior",
	int(pb.EDotaUserMessages_DOTA_UM_QoP_ArcanaSummary):          "mango.CDOTAUserMsg_QoP_ArcanaSummary",
	int(pb.EDotaUserMessages_DOTA_UM_HotPotato_Created):          "mango.CDOTAUserMsg_HotPotato_Created",
	int(pb.EDotaUserMessages_DOTA_UM_HotPotato_Exploded):         "mango.CDOTAUserMsg_HotPotato_Exploded",
	int(pb.EDotaUserMessages_DOTA_UM_WK_Arcana_Progress):         "mango.CDOTAUserMsg_WK_Arcana_Progress",
	int(pb.EDotaUserMessages_DOTA_UM_GuildChallenge_Progress):    "mango.CDOTAUserMsg_GuildChallenge_Progress",
	int(pb.EDotaUserMessages_DOTA_UM_WRArcanaProgress):           "mango.CDOTAUserMsg_WRArcanaProgress",
	int(pb.EDotaUserMessages_DOTA_UM_WRArcanaSummary):            "mango.CDOTAUserMsg_WRArcanaSummary",
	int(pb.EDotaUserMessages_DOTA_UM_ChatMessage):                "mango.CDOTAUserMsg_ChatMessage",
}

func GetEmbdeddedType(kind int) (string, proto.Message, error) {
	t, ok := EmbeddedTypes[kind]
	if !ok {
		return t, nil, fmt.Errorf("unknown embedded message type: %d", kind)
	}
	name := protoreflect.FullName(t)
	cls, err := protoregistry.GlobalTypes.FindMessageByName(name)
	if err != nil {
		return t, nil, err
	}
	data := cls.New().Interface()
	return t, data, nil
}

func RegisterHandler(kind int, handler EmbeddedHandler) {
	_, ok := EmbdeddedHandlers[kind]
	if !ok {
		EmbdeddedHandlers[kind] = []EmbeddedHandler{handler}
	} else {
		EmbdeddedHandlers[kind] = append(EmbdeddedHandlers[kind], handler)
	}
}
