CREATE OR REPLACE FUNCTION public.before_insert_user()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
	begin
		new.profileImg = '';
		new.discribtion = '';
		new.rating = 0;
		new.balance = 0; 
		new.isPremiumUser = false;
		new.isActiveUser = true;
		new.historyCount = 0;
		new.responsibility = 100;
		new.doneOnTime = 100;
		new.answerSpead = 100;
		new.registrationDate = NOW();
		return new;
	END;
$function$
;
create trigger bs_user before
insert
    on
    public.users for each row execute function before_insert_user()





CREATE OR REPLACE PROCEDURE public.changeactivity(IN userid integer)
 LANGUAGE plpgsql
AS $procedure$
	DECLARE
    	isActive boolean;
	begin
		select isActiveUser into isActive from Users where id = userId;
		if isActive then
			update Users set isActiveUser = false where id = userId;
		else
			update Users set isActiveUser = true where id = userId;
		end if;
	END;
$procedure$
;

CREATE OR REPLACE PROCEDURE public.changeactivityforoffer(IN userid integer, offerId integer)
 LANGUAGE plpgsql
AS $procedure$
	DECLARE
    	isActiveTmp boolean;
	begin
		select isActive into isActiveTmp from Offers where id = offerId and fkUserOwner = userId;
		if isActiveTmp then
			update Offers set isActive = false where id = offerId;
		else
			update Offers set isActive = true where id = offerId;
		end if;
	END;
$procedure$
;

CREATE OR REPLACE PROCEDURE public.makechat(IN userid integer, IN offerid integer, IN messagetextfromuser text)
 LANGUAGE plpgsql
AS $procedure$
	DECLARE
    	secondUserId integer;
    	tmpChatId integer;
    	tmpMessageId integer;
	begin
		select fkUserOwner into secondUserId from Offers where id = offerId;
		if (SELECT count(id) FROM chats where fkUserFirst = userid or fkUserSecond = userid) != 0 then
			select id into tmpChatId from chats where fkUserFirst = userid or fkUserSecond = userid;
			insert into messages (fkChatId, messageText, fkUserId, sendTime) values (tmpChatId, messagetextFromUser, userId, NOW()) RETURNING id into tmpMessageId;
			update chats set lastMessage = tmpMessageId where id = tmpChatId;
		else
			insert into chats (fkUserFirst, fkUserSecond) values (userId, secondUserId) RETURNING id into tmpChatId;
			insert into messages (fkChatId, messageText, fkUserId, sendTime) values (tmpChatId, messagetextFromUser, userId, NOW()) RETURNING id into tmpMessageId;
			update chats set lastMessage = tmpMessageId where id = tmpChatId;
		end if;
	END;
$procedure$
;



CREATE OR REPLACE FUNCTION public.getchatsbyuserid(_userid integer)
 RETURNS TABLE(id integer, fkusersecond integer, fkusername character varying, fkuserimg character varying, messagetext text, sendtime timestamp without time zone)
 LANGUAGE plpgsql
AS $function$
	begin
		RETURN QUERY SELECT t1.id, t1.fkUserSecond, t2.name, t2.profileImg, t3.messageText, t3.sendTime 
			FROM Chats as t1 LEFT JOIN Users as t2 on t1.fkUserSecond = t2.id LEFT JOIN Messages as t3
			on t1.lastMessage = t3.id WHERE t1.fkUserFirst = _userId union SELECT t1.id, t1.fkUserFirst, t2.name, t2.profileImg, t3.messageText, t3.sendTime 
			FROM Chats as t1 LEFT JOIN Users as t2 on t1.fkUserFirst = t2.id LEFT JOIN Messages as t3
			on t1.lastMessage = t3.id WHERE t1.fkUserSecond = _userId;
	END;
$function$
;

CREATE OR REPLACE PROCEDURE public.saveMessage(_chatId integer, _message varchar, _userId integer, _fileName varchar)
 LANGUAGE plpgsql
AS $procedure$
	DECLARE
    	tmpMessageId integer;
	begin
		insert into messages (fkChatId, messageText, fkUserId, sendTime, filepath) values (_chatId, _message, _userId, NOW(), _fileName) RETURNING id into tmpMessageId;
		update chats set lastMessage = tmpMessageId where id = _chatId;
	END;
$procedure$
;;
